package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
)

func ScanWorkSpaces(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult
	client := workspaces.NewFromConfig(cfg)

	input := &workspaces.DescribeWorkspacesInput{}
	for {
		resp, err := client.DescribeWorkspaces(ctx, input)
		if err != nil {
			return results
		}

		for _, ws := range resp.Workspaces {
			if ws.IpAddress != nil && *ws.IpAddress != "" {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      region,
					Service:     "WorkSpaces",
					ResourceID:  aws.ToString(ws.WorkspaceId),
					DNSName:     aws.ToString(ws.IpAddress),
					Extra:       aws.ToString(ws.UserName),
					PublicIP:    aws.ToString(ws.IpAddress),
				})
			}
		}

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		input.NextToken = resp.NextToken
	}

	return results
}

