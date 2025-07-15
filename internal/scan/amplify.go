package scan

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
)

func ScanAmplify(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := amplify.NewFromConfig(cfg)

	appsResp, err := client.ListApps(ctx, &amplify.ListAppsInput{})
	if err != nil {
		return results
	}

	for _, app := range appsResp.Apps {
		url := aws.ToString(app.DefaultDomain)
		if strings.Contains(url, ".amplifyapp.com") {
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "Amplify",
				ResourceID:  aws.ToString(app.Name),
				DNSName:     url,
				Extra:       "Amplify Hosting",
				PublicIP:    url,
			})
		}
	}

	return results
}

