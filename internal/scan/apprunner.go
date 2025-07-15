package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
)

func ScanAppRunner(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := apprunner.NewFromConfig(cfg)

	resp, err := client.ListServices(ctx, &apprunner.ListServicesInput{})
	if err != nil {
		return results
	}

	for _, svc := range resp.ServiceSummaryList {
		url := aws.ToString(svc.ServiceUrl)
		if url != "" {
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "AppRunner",
				ResourceID:  aws.ToString(svc.ServiceName),
				DNSName:     url,
				Extra:       string(svc.Status),
				PublicIP:    "",                     // ✅ Corrected
				ScanTarget:  PickScanTarget("", url), // ✅ Corrected
				URL:         url,                    // Optional
			})
		}
	}

	return results
}

