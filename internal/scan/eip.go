package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func ScanEIP(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult

	client := ec2.NewFromConfig(cfg)

	resp, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return results
	}

	for _, addr := range resp.Addresses {
		publicIP := aws.ToString(addr.PublicIp)
		if publicIP == "" {
			continue
		}

		allocID := aws.ToString(addr.AllocationId)
		instanceID := aws.ToString(addr.InstanceId)
		networkInterface := aws.ToString(addr.NetworkInterfaceId)

		extra := "unattached"
		if instanceID != "" {
			extra = "attached to instance " + instanceID
		} else if networkInterface != "" {
			extra = "attached to network interface " + networkInterface
		}

		results = append(results, ScanResult{
			AccountID:   accountID,
			AccountName: accountName,
			Region:      region,
			Service:     "eip",
			ResourceID:  allocID,
			DNSName:     "",
			PublicIP:    publicIP,
			Extra:       extra,
			ScanTarget:  publicIP, // ✅ IP used directly
			URL:         "",       // ✅ Do not populate
		})
	}

	return results
}

