package main

import (
  "os"
  "syscall"
  "time"

  log "github.com/sirupsen/logrus"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/ec2metadata"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
)

var start Time

func main() {
  start = time.Now()

  configureLogging("debug", "text")

  log.Debug("Starting AWS API session")
  ec2md_client := ec2metadata.New(session.New())

  log.Debug("Connecting to AWS EC2 metadata service")
  if ! ec2md_client.Available() {
    log.Fatal("Unable to access the metadata service, are you running this binary from an AWS EC2 instance?")
    exit(1)
  }

  log.Debug("Querying AZ from metadata service")
  az, err := ec2md_client.GetMetadata("placement/availability-zone")

  if err != nil {
    log.Fatal("Unable to retrieve the AZ from the EC2 metadata endpoint")
    exit(1)
  }

  log.Infof("Found AZ: '%s'", az )
  region := az[:len(az)-1]

  log.Infof("Computed region : '%s'", region )

  log.Debug("Querying instance-id from metadata service")
  id, err := ec2md_client.GetMetadata("instance-id")

  if err != nil {
    log.Fatal("Unable to retrieve the instance id from the EC2 metadata endpoint")
    exit(1)
  }

  log.Infof("Found instance-id : '%s'", id )

  log.Debug("Querying instance name tag from EC2 api endpoint")

  ec2_client := ec2.New(session.New(&aws.Config{
  	Region: aws.String(region),
  }))

  tags, err := ec2_client.DescribeTags(&ec2.DescribeTagsInput{
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
                  aws.String("Name"),
              },
          },
          {
              Name: aws.String("resource-id"),
              Values: []*string{
                  aws.String(id),
              },
          },
      },
  })

  if err != nil {
      if aerr, ok := err.(awserr.Error); ok {
          switch aerr.Code() {
          default:
              log.Fatal(aerr.Error())
              exit(1)
          }
      } else {
          log.Fatal(err.Error())
          exit(1)
      }
      return
  }

  if len(tags.Tags) != 1 {
      log.Fatal("Unable to fetch the EC2 instance name tag")
      exit(1)
  }

  if *tags.Tags[0].Key != "Name" {
    log.Fatalf("The tag fetched is not correct : '%s'", *tags.Tags[0].Key)
    exit(1)
  }

  nameTag := *tags.Tags[0].Value
  log.Infof("Found instance name tag : '%s'", nameTag)

  if nameTag[len(nameTag)-5:] == id[2:7] {
    log.Infof("Instance ID already found in the instance name/hostname : '%s', exiting..", *tags.Tags[0].Value)
    exit(0)
  }

  hostname := nameTag + "-" + id[2:7]

  log.Infof("Computed unique hostname : '%s'", hostname)

  log.Infof("Setting instance hostname locally")
  syscall.Sethostname([]byte(hostname))

  _, err = ec2_client.CreateTags(&ec2.CreateTagsInput{
      Resources: []*string{
        aws.String(id),
      },
      Tags: []*ec2.Tag{
          {
              Key:   aws.String("Name"),
              Value: aws.String(hostname),
          },
      },
  })

  if err != nil {
      if aerr, ok := err.(awserr.Error); ok {
          switch aerr.Code() {
          default:
              log.Fatal(aerr.Error())
              exit(1)
          }
      } else {
          log.Fatal(err.Error())
          exit(1)
      }
      return
  }

  exit(0)
}

func exit(code int) {
  log.Debugf("Executed in %s, exiting..", time.Since(start))
  os.Exit(code)
}
