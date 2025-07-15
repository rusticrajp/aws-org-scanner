package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
)

func ScanELB(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	client := elasticloadbalancing.NewFromConfig(cfg)
	out, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancing.DescribeLoadBalancersInput{})
	if err != nil || len(out.LoadBalancerDescriptions) == 0 {
		return nil
	}

	var results []ScanResult
	for _, lb := range out.LoadBalancerDescriptions {
		results = append(results, ScanResult{
			AccountID:   accountID,
			AccountName: accountName,
			Region:      region,
			Service:     "elb",
			ResourceID:  aws.ToString(lb.LoadBalancerName),
			DNSName:     aws.ToString(lb.DNSName),
			Extra:       fmt.Sprintf("%v", lb.ListenerDescriptions),
			ScanTarget:  aws.ToString(lb.DNSName),
		})
	}
	return results
}

