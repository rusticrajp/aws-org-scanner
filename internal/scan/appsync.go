package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

func ScanAppSync(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := appsync.NewFromConfig(cfg)

	resp, err := client.ListGraphqlApis(ctx, &appsync.ListGraphqlApisInput{})
	if err != nil {
		return results
	}

	for _, api := range resp.GraphqlApis {
		if api.Uris != nil {
			if graphqlURL, ok := api.Uris["GRAPHQL"]; ok && graphqlURL != "" {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      region,
					Service:     "AppSync",
					ResourceID:  aws.ToString(api.Name),
					DNSName:     graphqlURL,
					Extra:       "GraphQL endpoint",
					PublicIP:    "",                         // ✅ fixed
					ScanTarget:  PickScanTarget("", graphqlURL), // ✅ fixed
					URL:         graphqlURL,                // ✅ optional
				})
			}
		}
	}

	return results
}

