package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func ScanEC2(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult

	client := ec2.NewFromConfig(cfg)

	resp, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return results
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			publicIP := aws.ToString(instance.PublicIpAddress)
			if publicIP == "" {
				continue // skip instances without public IP
			}

			instanceID := aws.ToString(instance.InstanceId)
			nameTag := ""
			for _, tag := range instance.Tags {
				if aws.ToString(tag.Key) == "Name" {
					nameTag = aws.ToString(tag.Value)
					break
				}
			}

			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "ec2",
				ResourceID:  instanceID,
				DNSName:     "",
				PublicIP:    publicIP,
				Extra:       nameTag,
				ScanTarget:  publicIP, // ✅ use public IP directly
				URL:         "",       // ✅ do not populate
			})
		}
	}

	return results
}

