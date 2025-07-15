package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func ScanLambda(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := lambda.NewFromConfig(cfg)

	// Paginate through all Lambda functions
	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			break
		}

		for _, fn := range page.Functions {
			fnName := aws.ToString(fn.FunctionName)

			// Check if function has a Function URL
			urlResp, err := client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
				FunctionName: &fnName,
			})
			if err != nil {
				continue // skip if no URL or permission denied
			}

			url := aws.ToString(urlResp.FunctionUrl)
			if url == "" {
				continue
			}

			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "lambda",
				ResourceID:  fnName,
				DNSName:     "",            // URL only, not raw DNS
				PublicIP:    "",            // not applicable
				Extra:       "function URL",
				ScanTarget:  url,           // ✅ set as scan target
				URL:         url,           // ✅ visible URL
			})
		}
	}

	return results
}

