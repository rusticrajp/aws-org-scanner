package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func ScanELBV2(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	client := elasticloadbalancingv2.NewFromConfig(cfg)
	out, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil || len(out.LoadBalancers) == 0 {
		return nil
	}

	var results []ScanResult
	for _, lb := range out.LoadBalancers {
		results = append(results, ScanResult{
			AccountID:   accountID,
			AccountName: accountName,
			Region:      region,
			Service:     "elbv2",
			ResourceID:  aws.ToString(lb.LoadBalancerName),
			DNSName:     aws.ToString(lb.DNSName),
			Extra:       fmt.Sprintf("Type: %s, Scheme: %s", lb.Type, lb.Scheme),
			ScanTarget:  aws.ToString(lb.DNSName),
		})
	}
	return results
}

