package scan

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)


// LoadAWSConfig loads the default AWS configuration
func LoadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
}

// GetOrgAccounts lists all active accounts in the AWS Organization
func GetOrgAccounts(ctx context.Context, cfg aws.Config) ([]Account, error) {
	orgClient := organizations.NewFromConfig(cfg)
	var accounts []Account

	paginator := organizations.NewListAccountsPaginator(orgClient, &organizations.ListAccountsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, acct := range page.Accounts {
			if acct.Status == "ACTIVE" && acct.Id != nil && acct.Name != nil {
				accounts = append(accounts, Account{
					ID:   *acct.Id,
					Name: *acct.Name,
				})
			}
		}
	}
	return accounts, nil
}

// MustGetCurrentAccountID returns the account ID from the caller identity
func MustGetCurrentAccountID(ctx context.Context, cfg aws.Config) string {
	stsClient := sts.NewFromConfig(cfg)
	out, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("❌ Failed to get caller identity: %v", err)
	}
	return *out.Account
}

// AssumeRole assumes an IAM role into the specified account
func AssumeRole(ctx context.Context, cfg aws.Config, accountID, roleName string) (aws.Config, error) {
	stsClient := sts.NewFromConfig(cfg)
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)

	out, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("ScanSession-" + accountID),
		DurationSeconds: aws.Int32(900),
	})
	if err != nil {
		return aws.Config{}, err
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		*out.Credentials.AccessKeyId,
		*out.Credentials.SecretAccessKey,
		*out.Credentials.SessionToken,
	))

	return config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(creds))
}

// GetEnabledRegions returns a list of all enabled AWS regions
func GetEnabledRegions(ctx context.Context, cfg aws.Config) []string {
	client := ec2.NewFromConfig(cfg)
	resp, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return []string{"us-east-1"} // fallback
	}

	var regions []string
	for _, r := range resp.Regions {
		if r.RegionName != nil {
			regions = append(regions, *r.RegionName)
		}
	}
	return regions
}

// PickScanTarget returns the best scan target: Public IP if available, otherwise DNS name.
func PickScanTarget(ip, dns string) string {
	if ip != "" {
		return ip
	}
	return dns
}

