package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

func ScanOpenSearch(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := opensearch.NewFromConfig(cfg)

	domainsResp, err := client.ListDomainNames(ctx, &opensearch.ListDomainNamesInput{})
	if err != nil {
		return results
	}

	for _, domain := range domainsResp.DomainNames {
		descResp, err := client.DescribeDomain(ctx, &opensearch.DescribeDomainInput{
			DomainName: domain.DomainName,
		})
		if err != nil {
			continue
		}

		status := descResp.DomainStatus
		if status == nil || status.Endpoint == nil {
			continue
		}

		// If VPCOptions is nil, the endpoint is public
		if status.VPCOptions == nil {
			endpoint := aws.ToString(status.Endpoint)
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "OpenSearch",
				ResourceID:  aws.ToString(domain.DomainName),
				DNSName:     endpoint,
				PublicIP:    "",                        // ✅ leave empty
				Extra:       "Public Domain",
				ScanTarget:  PickScanTarget("", endpoint), // ✅ fixed
				URL:         "",                        // Optional: could add AWS console URL
			})
		}
	}

	return results
}

