package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

func ScanRedshift(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := redshift.NewFromConfig(cfg)

	resp, err := client.DescribeClusters(ctx, &redshift.DescribeClustersInput{})
	if err != nil {
		return results
	}

	for _, cluster := range resp.Clusters {
		if cluster.PubliclyAccessible != nil && *cluster.PubliclyAccessible {
			endpoint := ""
			if cluster.Endpoint != nil {
				endpoint = aws.ToString(cluster.Endpoint.Address)
			}

			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "redshift",
				ResourceID:  aws.ToString(cluster.ClusterIdentifier),
				DNSName:     endpoint,
				PublicIP:    "",
				Extra:       "publicly accessible",
				ScanTarget:  PickScanTarget("", endpoint),
				URL: "", 

			})
		}
	}

	return results
}

