package awscmd

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type InputEc2List struct {
	Region string
}

type OutputEc2List struct {
	Region    string
	Instances []*ec2Instance
}

type ec2Instance struct {
	ID          string
	Name        string
	IPv4        string
	IPv4private string
	Type        string
	State       string
	AMI         string
	Zone        string
	LaunchTime  time.Time
}

func Ec2List(ctx context.Context, input *InputEc2List, w io.Writer) (*OutputEc2List, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	response, err := svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	ec2s := []*ec2Instance{}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			name := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}
			row := &ec2Instance{
				ID:          toS(instance.InstanceId),
				Name:        name,
				IPv4:        toS(instance.PublicIpAddress),
				IPv4private: toS(instance.PrivateIpAddress),
				Type:        toS(instance.InstanceType),
				State:       ec2InstanceStatus(*instance.State.Code),
				AMI:         toS(instance.ImageId),
				Zone:        toS(instance.Placement.AvailabilityZone),
				LaunchTime:  *instance.LaunchTime,
			}
			ec2s = append(ec2s, row)
		}
	}

	return &OutputEc2List{
		Region:    input.Region,
		Instances: ec2s,
	}, nil
}

func toS(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ec2InstanceStatus(state int64) string {
	switch state {
	case 16:
		return "running"
	case 32:
		return "shutting-down"
	case 48:
		return "terminated"
	case 64:
		return "stopping"
	case 80:
		return "stopped"
	default: // 0
		return "pending"
	}
}
