package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ScanS3(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult

	client := s3.NewFromConfig(cfg)

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return results
	}

	for _, bucket := range output.Buckets {
		bucketName := aws.ToString(bucket.Name)

		// Try to get bucket location
		locResp, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: &bucketName,
		})
		bucketRegion := region
		if err == nil && locResp.LocationConstraint != "" {
			bucketRegion = string(locResp.LocationConstraint)
		}

		// Check if bucket has public access disabled
		policyStatus, err := client.GetBucketPolicyStatus(ctx, &s3.GetBucketPolicyStatusInput{
			Bucket: &bucketName,
		})

		if err == nil && policyStatus.PolicyStatus != nil && policyStatus.PolicyStatus.IsPublic != nil && *policyStatus.PolicyStatus.IsPublic {
			results = append(results, ScanResult{
				AccountID:   accountID,
				AccountName: accountName,
				Region:      bucketRegion,
				Service:     "s3",
				ResourceID:  bucketName,
				DNSName:     fmt.Sprintf("%s.s3.amazonaws.com", bucketName),
				PublicIP:    "",
				Extra:       "Bucket is public",
				ScanTarget:  bucketName,
				URL:         fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName),
			})
		}
	}

	return results
}

