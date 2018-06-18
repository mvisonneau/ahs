package main

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time

func run(c *cli.Context) error {
	start = time.Now()
	configureLogging(cfg.Log.Level, cfg.Log.Format)

	log.Debug("Starting AWS MDS API session")
	mdsClient, err := getAWSMDSClient()
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	log.Debug("Fetching current instance-id from MDS API")
	instanceID, err := getInstanceID(mdsClient)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}
	log.Infof("Found instance-id : '%s'", id)

	log.Debug("Fetching current AZ from MDS API")
	az, err := getInstanceAZ(mdsClient)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}
	log.Infof("Found AZ: '%s'", az)

	region := computeRegionFromAZ(az)
	log.Infof("Computed region : '%s'", region)

	log.Debug("Starting AWS EC2 API session")
	ec2Client := getAWSEC2Client(region)

	log.Debugf("Querying Input Tag '%s' from EC2 API", cfg.InputTag)
	inputTagValue, err := getInputTagValue(cfg.InputTag, instanceID, ec2Client)
	if err != nil {
		return exit(cli.NewExitError(analyzeEC2APIErrors(err), 1))
	}
	log.Infof("Found instance name tag : '%s'", inputTagValue)

	if inputTagValue[len(inputTagValue)-cfg.IDLength:] == instanceID[2:2+cfg.IDLength] {
		hostname := inputTagValue
		log.Infof("Instance ID already found in the instance tag : '%s', reusing this value", inputTagValue)
	} else {
		hostname := computeHostname(inputTagValue, cfg.Separator, instanceID, cfg.IDLength)
		log.Infof("Computed unique hostname : '%s'", hostname)
	}

	if !cfg.DryRun {
		log.Infof("Setting instance hostname locally")
		err = setHostname(hostname)
		if err != nil {
			return exit(cli.NewExitError(err.Error(), 1))
		}

		log.Infof("Setting hostname on configured instance tag '%s'", cfg.OutputTag)
		err = setOutputTagValue(cfg.OutputTag, hostname, instanceID, ec2Client)
		if err != nil {
			return exit(cli.NewExitError(analyzeEC2APIErrors(err), 1))
		}
	} else {
		log.Infof("Setting instance hostname locally (dry-run)")
		log.Infof("Setting hostname on configured instance tag '%s' (dry-run)", cfg.OutputTag)
	}

	return exit(nil)
}

func getAWSMDSClient() (*ec2metadata.EC2Metadata, error) {
	client := ec2metadata.New(session.New())

	if !client.Available() {
		return client, errors.New("Unable to access the metadata service, are you running this binary from an AWS EC2 instance?")
	}

	return client, nil
}

func getAWSEC2Client(region string) (client *ec2.EC2) {
	client = ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	return
}

func getInstanceAZ(c *ec2metadata.EC2Metadata) (az string, err error) {
	az, err = c.GetMetadata("placement/availability-zone")
	return
}

func computeRegionFromAZ(az string) string {
	return az[:len(az)-1]
}

func getInstanceID(c *ec2metadata.EC2Metadata) (id string, err error) {
	id, err = c.GetMetadata("instance-id")
	return
}

func getInputTagValue(tag string, instanceID string, c *ec2.EC2) (string, error) {
	log.Debug("Querying instance name tag from EC2 api endpoint")

	tags, err := c.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String("instance"),
				},
			},
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String(tag),
				},
			},
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
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

	if *tags.Tags[0].Key != tag {
		return "", fmt.Errorf("The tag fetched is not correct : '%s'", *tags.Tags[0].Key)
	}

	return *tags.Tags[0].Value, nil
}

func analyzeEC2APIErrors(err error) string {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr.Error()
		}
		return err.Error()
	}
	return ""
}

func computeHostname(base string, separator string, id string, IDLength int) string {
	return base + separator + id[2:2+IDLength]
}

func setHostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}

func setOutputTagValue(tag string, hostname string, instanceID string, c *ec2.EC2) (err error) {
	_, err = c.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(instanceID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tag),
				Value: aws.String(hostname),
			},
		},
	})

	return
}

func exit(err error) error {
	log.Debugf("Executed in %s, exiting..", time.Since(start))
	return err
}
