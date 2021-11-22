package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
	"github.com/urfave/cli/v2"
)

const (
	sequential = "sequential"
)

// Params of the app
type Params struct {
	Backoff   *backoff.Backoff
	InputTag  string
	OutputTag string
	Separator string
}

// Clients of AWS libs
type Clients struct {
	Autoscaling *autoscaling.AutoScaling
	EC2         *ec2.EC2
	MDS         *ec2metadata.EC2Metadata
}

// Values computed/generated
type Values struct {
	AZ           string
	Base         string
	Hostname     string
	InstanceID   string
	Region       string
	SequentialID int
}

// Run is the main handler for all our functions
// TODO: Break it apart in smaller ones
func Run(ctx *cli.Context) (int, error) {
	if err := configure(ctx); err != nil {
		return 1, err
	}

	if user, err := user.Current(); err != nil {
		return 1, fmt.Errorf("Unable to determine current user")
	} else if user.Username != "root" {
		return 1, fmt.Errorf("You have to run this function as root")
	}

	p := &Params{
		Backoff: &backoff.Backoff{
			Min:    100 * time.Millisecond,
			Max:    120 * time.Second,
			Factor: 2,
			Jitter: false,
		},
		InputTag:  ctx.String("input-tag"),
		OutputTag: ctx.String("output-tag"),
		Separator: ctx.String("separator"),
	}

	c := &Clients{}
	v := &Values{}

	// Configure MDS Client
	if err := c.getAWSMDSClient(); err != nil {
		return 1, err
	}

	// Fetch current AZ
	var err error
	v.AZ, err = c.getInstanceAZ()
	if err != nil {
		return 1, err
	}

	// Compute region from AZ
	v.Region, err = computeRegionFromAZ(v.AZ)
	if err != nil {
		return 1, err
	}

	// Configure EC2 Client
	if err := c.getAWSEC2Client(v.Region); err != nil {
		return 1, err
	}

	// Fetch instance ID
	v.InstanceID, err = c.getInstanceID()
	if err != nil {
		return 1, err
	}

	// Fetch the value of the input-tag and use it a base for the hostname
	for {
		v.Base, err = c.getBaseFromInputTag(p.InputTag, v.InstanceID)
		if err != nil {
			d := p.Backoff.Duration()
			if d == 60*time.Second {
				return 1, err
			}
			log.Infof("%s, retrying in %s", err, d)
			time.Sleep(d)
		} else {
			p.Backoff.Reset()
			break
		}
	}

	switch ctx.Command.FullName() {
	case "instance-id":
		v.Hostname, err = computeHostnameWithInstanceID(v.Base, p.Separator, v.InstanceID, ctx.Int("length"))
	case sequential:
		var instanceGroup string
		instanceGroup, err = c.findInstanceGroupTagValue(ctx.String("instance-group-tag"), v.InstanceID)
		if err != nil {
			return 1, err
		}

		if !ctx.Bool("respect-azs") {
			v.SequentialID, err = c.findAvailableSequentialIDPerRegion(instanceGroup, ctx.String("instance-group-tag"), ctx.String("instance-sequential-id-tag"))
			if err != nil {
				return 1, err
			}
		} else {
			// Configure Autoscaling Client
			if err = c.getAWSAutoscalingClient(v.Region); err != nil {
				return 1, err
			}

			v.SequentialID, err = c.findAvailableSequentialIDPerAZ(v.AZ, instanceGroup, ctx.String("instance-group-tag"), ctx.String("instance-sequential-id-tag"))
			if err != nil {
				return 1, err
			}
		}
		v.Hostname, err = computeSequentialHostname(v.Base, p.Separator, v.SequentialID)
	default:
		return 1, fmt.Errorf("Function %v is not implemented", ctx.Command.FullName())
	}

	if err != nil {
		return 1, err
	}

	if !ctx.Bool("dry-run") {
		log.Infof("Setting instance hostname locally")
		if err := setSystemHostname(v.Hostname); err != nil {
			return 1, err
		}

		if ctx.Bool("persist-hostname") {
			log.Infof("Persist hostname in /etc/hostname")
			if err := updateHostnameFile(v.Hostname); err != nil {
				return 1, err
			}
		}

		if ctx.Bool("persist-hosts") {
			log.Infof("Persist hostname in /etc/hosts")
			if err := updateHostsFile(v.Hostname); err != nil {
				return 1, err
			}
		}

		log.Infof("Setting hostname on configured instance output tag '%s'", p.OutputTag)
		if err := c.setTagValue(v.InstanceID, p.OutputTag, v.Hostname); err != nil {
			return 1, err
		}

		if ctx.Command.FullName() == sequential {
			log.Infof("Setting instance sequential id (%d) on configured tag '%s'", v.SequentialID, ctx.String("instance-sequential-id-tag"))
			if err := c.setTagValue(v.InstanceID, ctx.String("instance-sequential-id-tag"), strconv.Itoa(v.SequentialID)); err != nil {
				return 1, err
			}
		}
	} else {
		log.Infof("Setting instance hostname locally (dry-run)")
		log.Infof("Setting hostname on configured instance tag '%s' (dry-run)", p.OutputTag)
		if ctx.Command.FullName() == sequential {
			log.Infof("Setting instance sequential id (%d) on configured tag '%s' (dry-run)", v.SequentialID, ctx.String("instance-sequential-id-tag"))
		}
	}

	return 0, nil
}

func (c *Clients) getAWSMDSClient() error {
	log.Debug("Starting AWS MDS API session")
	c.MDS = ec2metadata.New(session.New())

	if !c.MDS.Available() {
		return errors.New("Unable to access the metadata service, are you running this binary from an AWS EC2 instance?")
	}

	return nil
}

func (c *Clients) getAWSEC2Client(region string) (err error) {
	re := regexp.MustCompile("[a-z]{2}-[a-z]+-\\d")
	if !re.MatchString(region) {
		return fmt.Errorf("Cannot start AWS EC2 client session with invalid region '%s'", region)
	}

	log.Debug("Starting AWS EC2 Client session")
	c.EC2 = ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	return
}

func (c *Clients) getAWSAutoscalingClient(region string) (err error) {
	re := regexp.MustCompile("[a-z]{2}-[a-z]+-\\d")
	if !re.MatchString(region) {
		return fmt.Errorf("Cannot start AWS Autoscaling client session with invalid region '%s'", region)
	}

	log.Debug("Starting AWS EC2 Client session")
	c.Autoscaling = autoscaling.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	return
}

func (c *Clients) getInstanceAZ() (az string, err error) {
	log.Debug("Fetching current AZ from MDS API")
	az, err = c.MDS.GetMetadata("placement/availability-zone")
	log.Infof("Found AZ: '%s'", az)
	return
}

func computeRegionFromAZ(az string) (region string, err error) {
	re := regexp.MustCompile("[a-z]{2}-[a-z]+-\\d[a-z]")
	if !re.MatchString(az) {
		err = fmt.Errorf("Cannot compute region from invalid availability-zone '%s'", az)
		return
	}

	region = az[:len(az)-1]
	log.Infof("Computed region : '%s'", region)
	return
}

func (c *Clients) getInstanceID() (iid string, err error) {
	log.Debug("Fetching current instance-id from MDS API")
	iid, err = c.MDS.GetMetadata("instance-id")
	log.Infof("Found instance-id : '%s'", iid)
	return
}

func (c *Clients) getBaseFromInputTag(inputTag, instanceID string) (string, error) {
	log.Infof("Querying input-tag '%s' from EC2 API", inputTag)
	instances, err := c.EC2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == inputTag {
					log.Debugf("Found input-tag '%s' : '%s' ", inputTag, *tag.Value)
					return *tag.Value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Instance doesn't contain input-tag '%s'", inputTag)
}

func getSystemHostname() (string, error) {
	return os.Hostname()
}

func (c *Clients) setTagValue(instanceID, tag, value string) (err error) {
	_, err = c.EC2.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(instanceID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tag),
				Value: aws.String(value),
			},
		},
	})

	return
}

func computeHostnameWithInstanceID(base, separator, instanceID string, length int) (string, error) {
	log.Info("Computing hostname with truncated instance-id")

	// remove i-
	awsInstanceID := instanceID[2:]
	truncatedID := truncateString(awsInstanceID, length)

	if base == truncatedID {
		log.Infof("Instance ID already found in the instance tag : '%s', reusing this value", base)
		return base, nil
	}

	splitHost := strings.Split(base, separator)
	log.Infof("splitHost : '%v'", splitHost)
	splitIndex := len(splitHost) - 1
	if len(splitHost) <= 1 {
		splitIndex = 1
	}
	hostnameIncluded := strings.Contains(awsInstanceID, splitHost[len(splitHost)-1])
	if !hostnameIncluded {
		splitIndex = len(splitHost)
	}

	hostnamePrefix := strings.Join(splitHost[:splitIndex], separator)
	hostname := strings.Join([]string{hostnamePrefix, truncatedID}, separator)

	return hostname, nil
}

func truncateString(str string, length int) string {
	if length <= 0 {
		return str
	}

	clampedLength := len(str)
	if length < clampedLength {
		clampedLength = length
	}

	return str[:clampedLength]
}

func computeSequentialHostname(base, separator string, sequentialID int) (string, error) {
	log.Info("Computing a hostname with sequential naming")

	re := regexp.MustCompile(".*-(\\d+)$")
	if re.MatchString(base) {
		log.Infof("Current input tag value already matches '.*-\\d+$', keeping '%s' as hostname", base)
		return base, nil
	}

	hostname := base + separator + strconv.Itoa(sequentialID)
	log.Infof("Computed unique hostname : '%s'", hostname)
	return hostname, nil
}

func (c *Clients) findInstanceGroupTagValue(groupTag, instanceID string) (string, error) {
	log.Debugf("Looking up the value of the tag '%s' of the instance", groupTag)
	tags, err := c.EC2.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String("instance"),
				},
			},
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String(groupTag),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	if len(tags.Tags) != 1 {
		return "", fmt.Errorf("Unexpected amount of tags retrieved : '%d',  expected 1", len(tags.Tags))
	}

	log.Debugf("Found instance-group value : '%s'", *tags.Tags[0].Value)
	return *tags.Tags[0].Value, nil
}

func (c *Clients) getASG(asgName string) (*autoscaling.Group, error) {
	log.Debugf("Looking for ASG '%s'", asgName)
	asgs, err := c.Autoscaling.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	})
	if err != nil {
		return nil, err
	}

	if len(asgs.AutoScalingGroups) != 1 {
		return nil, fmt.Errorf("Unexpected amount of asgs retrieved : '%d',  expected 1", len(asgs.AutoScalingGroups))
	}

	return asgs.AutoScalingGroups[0], nil
}

func (c *Clients) getASGAZs(asgName string) ([]*string, error) {
	asg, err := c.getASG(asgName)
	if err != nil {
		return nil, err
	}

	log.Debugf("Found '%d' AZ(s)", len(asg.AvailabilityZones))
	return asg.AvailabilityZones, nil
}

func (c *Clients) getASGMaxInstances(asgName string) (int, error) {
	log.Debugf("Getting maximum size of the ASG '%s'", asgName)
	asg, err := c.getASG(asgName)
	if err != nil {
		return 0, err
	}

	log.Debugf("Found ASG '%s' max size : %d", asgName, int(*asg.MaxSize))
	return int(*asg.MaxSize), nil
}

func (c *Clients) findAvailableSequentialIDPerRegion(instanceGroup, groupTag, sequentialIDTag string) (int, error) {
	log.Debugf("Looking up instances that belong to the same group within the region")
	instances, err := c.EC2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + groupTag),
				Values: []*string{
					aws.String(instanceGroup),
				},
			},
		},
	})
	if err != nil {
		return -1, err
	}

	return computeMostAdequateSequentialID(instances, sequentialIDTag, 1, 1)
}

func (c *Clients) findAvailableSequentialIDPerAZ(instanceAZ, instanceGroup, groupTag, sequentialIDTag string) (int, error) {
	log.Debugf("Looking up how many AZs are configured on the ASG")
	azs, err := c.getASGAZs(instanceGroup)
	if err != nil {
		return -1, err
	}

	max, err := c.getASGMaxInstances(instanceGroup)
	if err != nil {
		return -1, err
	}

	log.Debugf("Looking up instances that belong to the same group within the AZ (%s)", instanceAZ)
	instances, err := c.EC2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(instanceAZ),
				},
			},
			{
				Name: aws.String("tag:" + groupTag),
				Values: []*string{
					aws.String(instanceGroup),
				},
			},
		},
	})
	if err != nil {
		return -1, err
	}

	// Get an offset based on the letter of the AZ
	var offset int
	azList := []string{"a", "b", "c", "d", "e", "f"}
	for i := range azList {
		if instanceAZ[len(instanceAZ)-1:] == azList[i] {
			offset = i + 1
			break
		}
	}

	computedID, err := computeMostAdequateSequentialID(instances, sequentialIDTag, offset, len(azs))
	if err != nil {
		return -1, err
	}

	if computedID > max {
		return -1, fmt.Errorf("Computed ID %d is higher than the size of the ASG.. (%d)", computedID, max)
	}

	return computedID, nil
}

func computeMostAdequateSequentialID(instances *ec2.DescribeInstancesOutput, sequentialIDTag string, offset, modulo int) (int, error) {
	var used []int
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			if *instance.State.Name == "running" {
				for _, tag := range instance.Tags {
					if *tag.Key == sequentialIDTag {
						v, err := strconv.Atoi(*tag.Value)
						if err != nil {
							return -1, err
						}

						skip := false
						for i := 0; i < len(used); i++ {
							if used[i] == v {
								log.Warnf("Found another running instance '%s' with sequential id '%d'!, skipping it for the count", *instance.InstanceId, v)
								skip = true
							}
						}

						if !skip {
							used = append(used, v)
							log.Debugf("Found running instance '%s' with sequential id '%d' ", *instance.InstanceId, v)
						}
					}
				}
			}
		}
	}

	if len(used) > 0 {
		sort.Ints(used)

		// if the instance holding the first id has disappeared, we get it
		if used[0] != offset {
			return offset, nil
		}

		// search if there are no other missed ids
		for i := 1; i < len(used); i++ {
			if used[i] != (i*modulo)+offset {
				return (i * modulo) + offset, nil
			}
		}

		// return an incremental one
		return (len(used) * modulo) + offset, nil
	}

	// if there is not a single instance, we start with the offset
	return offset, nil
}

func updateHostnameFile(hostname string) error {
	return ioutil.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0o644)
}

func updateHostsFile(hostname string) error {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	hosts.RemoveHosts([]string{hostname})
	hosts.AddHost("127.0.0.1", hostname)

	return hosts.Save()
}
