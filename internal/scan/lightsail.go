package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
)

func ScanLightsail(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := lightsail.NewFromConfig(cfg)

	// Instances
	instResp, err := client.GetInstances(ctx, &lightsail.GetInstancesInput{})
	if err == nil {
		for _, inst := range instResp.Instances {
			if inst.PublicIpAddress != nil && *inst.PublicIpAddress != "" {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      region,
					Service:     "Lightsail",
					ResourceID:  aws.ToString(inst.Name),
					DNSName:     *inst.PublicIpAddress,
					Extra:       "Instance with Public IP",
					PublicIP:    *inst.PublicIpAddress,
				})
			}
		}
	}

	// Static IPs
	ipResp, err := client.GetStaticIps(ctx, &lightsail.GetStaticIpsInput{})
	if err == nil {
		for _, ip := range ipResp.StaticIps {
			if ip.IpAddress != nil && *ip.IpAddress != "" {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      region,
					Service:     "Lightsail",
					ResourceID:  aws.ToString(ip.Name),
					DNSName:     *ip.IpAddress,
					Extra:       "Static IP",
					PublicIP:    *ip.IpAddress,
				})
			}
		}
	}

	return results
}

