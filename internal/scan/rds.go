package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func ScanRDS(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := rds.NewFromConfig(cfg)

	resp, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return results
	}

	for _, db := range resp.DBInstances {
		if db.PubliclyAccessible != nil && *db.PubliclyAccessible {
			endpoint := ""
			if db.Endpoint != nil {
				endpoint = aws.ToString(db.Endpoint.Address)
			}

			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "rds",
				ResourceID:  aws.ToString(db.DBInstanceIdentifier),
				DNSName:     endpoint,
				PublicIP:    "",
				Extra:       "publicly accessible",
				ScanTarget:  PickScanTarget("", endpoint), // ✅ FIXED
				URL: fmt.Sprintf("https://%s.console.aws.amazon.com/rds/home?region=%s#database:id=%s;is-cluster=false",
					region, region, aws.ToString(db.DBInstanceIdentifier)),
			})
		}
	}

	return results
}

