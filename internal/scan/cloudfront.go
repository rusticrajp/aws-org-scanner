package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

func ScanCloudFront(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := cloudfront.NewFromConfig(cfg)

	resp, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{})
	if err != nil {
		return results
	}

	for _, dist := range resp.DistributionList.Items {
		// If aliases exist, record those
		if len(dist.Aliases.Items) > 0 {
			for _, alias := range dist.Aliases.Items {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      "global",
					Service:     "cloudfront",
					ResourceID:  aws.ToString(dist.Id),
					DNSName:     alias,
					PublicIP:    "",
					Extra:       "alias",
					ScanTarget:  alias,
					URL:         fmt.Sprintf("https://%s", alias),
				})
			}
		} else {
			// No alias — use the default domain
			defaultDomain := aws.ToString(dist.DomainName)
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      "global",
				Service:     "cloudfront",
				ResourceID:  aws.ToString(dist.Id),
				DNSName:     defaultDomain,
				PublicIP:    "",
				Extra:       "default domain",
				ScanTarget:  defaultDomain,
				URL:         fmt.Sprintf("https://%s", defaultDomain),
			})
		}
	}

	return results
}

