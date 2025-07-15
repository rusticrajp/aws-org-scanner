package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

func ScanAPIGateway(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	client := apigateway.NewFromConfig(cfg)
	out, err := client.GetRestApis(ctx, &apigateway.GetRestApisInput{})
	if err != nil || len(out.Items) == 0 {
		return nil
	}

	var results []ScanResult
	for _, api := range out.Items {
		id := aws.ToString(api.Id)
		name := aws.ToString(api.Name)
		url := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", id, region)

		results = append(results, ScanResult{
			AccountID:   accountID,
			AccountName: accountName,
			Region:      region,
			Service:     "apigateway",
			ResourceID:  id,
			Extra:       name,
			URL:         url,
			ScanTarget:  url,
		})
	}
	return results
}

