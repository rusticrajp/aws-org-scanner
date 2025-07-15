package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func ScanRoute53(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult

	client := route53.NewFromConfig(cfg)

	resp, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		return results
	}

	for _, zone := range resp.HostedZones {
		zoneID := aws.ToString(zone.Id)
		zoneName := aws.ToString(zone.Name)

		// Skip private zones
		if zone.Config != nil && zone.Config.PrivateZone {
			continue
		}

		results = append(results, ScanResult{
			AccountID:   accountID,
			AccountName: accountName,
			Region:      "global",
			Service:     "route53",
			ResourceID:  zoneID,
			DNSName:     zoneName,
			PublicIP:    "",
			Extra:       "public hosted zone",
			ScanTarget:  PickScanTarget("", zoneName),
			URL: "",
		})
	}

	return results
}

