package scan

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
)

func ScanElasticBeanstalk(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := elasticbeanstalk.NewFromConfig(cfg)

	apps, err := client.DescribeEnvironments(ctx, &elasticbeanstalk.DescribeEnvironmentsInput{
		IncludeDeleted: aws.Bool(false),
	})
	if err != nil {
		return results
	}

	for _, env := range apps.Environments {
		url := aws.ToString(env.CNAME)
		if url != "" && strings.HasSuffix(url, ".elasticbeanstalk.com") {
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      region,
				Service:     "ElasticBeanstalk",
				ResourceID:  aws.ToString(env.EnvironmentName),
				DNSName:     url,
				Extra:       string(env.Status),
				PublicIP:    url,
			})
		}
	}

	return results
}

