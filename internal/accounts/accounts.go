package accounts

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

type AccountInfo struct {
	ID     string
	Name   string
	IsRoot bool
}

const RoleNameToAssume = "OrganizationAccountAccessRole"

func FetchAllAccounts(ctx context.Context, rootCfg aws.Config) []AccountInfo {
	orgClient := organizations.NewFromConfig(rootCfg)

	orgDetail, err := orgClient.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		log.Fatalf("❌ Failed to describe organization: %v", err)
	}
	managementAccountID := *orgDetail.Organization.MasterAccountId

	accountList, err := orgClient.ListAccounts(ctx, &organizations.ListAccountsInput{})
	if err != nil {
		log.Fatalf("❌ Failed to list AWS accounts: %v", err)
	}

	var result []AccountInfo
	for _, acct := range accountList.Accounts {
		if acct.Status != "ACTIVE" || acct.Id == nil {
			continue
		}
		result = append(result, AccountInfo{
			ID:     *acct.Id,
			Name:   *acct.Name,
			IsRoot: *acct.Id == managementAccountID,
		})
	}
	return result
}

func AssumeOrUseRoot(ctx context.Context, rootCfg aws.Config, acct AccountInfo) aws.Config {
	if acct.IsRoot {
		return rootCfg
	}

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", acct.ID, RoleNameToAssume)
	stsClient := sts.NewFromConfig(rootCfg)

	resp, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("OrgScanSession-" + acct.ID),
		DurationSeconds: aws.Int32(900),
	})
	if err != nil {
		log.Printf("[!] Failed to assume role for account %s (%s): %v\n", acct.ID, acct.Name, err)
		return aws.Config{}
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		*resp.Credentials.AccessKeyId,
		*resp.Credentials.SecretAccessKey,
		*resp.Credentials.SessionToken,
	))

	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(creds))
	if err != nil {
		log.Printf("[!] Failed to load config for assumed role: %v", err)
	}
	return cfg
}

func GetEnabledRegions(ctx context.Context, cfg aws.Config) []string {
	client := ec2.NewFromConfig(cfg)
	resp, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil
	}
	var regions []string
	for _, r := range resp.Regions {
		if r.RegionName != nil {
			regions = append(regions, *r.RegionName)
		}
	}
	return regions
}

