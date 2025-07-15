package scan

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
)

func ScanGlobalAccelerator(ctx context.Context, cfg aws.Config, accountID, accountName, region string) []ScanResult {
	var results []ScanResult

	client := globalaccelerator.NewFromConfig(cfg)

	accelerators, err := client.ListAccelerators(ctx, &globalaccelerator.ListAcceleratorsInput{})
	if err != nil {
		return results
	}

	for _, acc := range accelerators.Accelerators {
		accArn := aws.ToString(acc.AcceleratorArn)
		accName := aws.ToString(acc.Name)
		ipSetResp, err := client.DescribeAccelerator(ctx, &globalaccelerator.DescribeAcceleratorInput{
			AcceleratorArn: &accArn,
		})
		if err != nil {
			continue
		}

		for _, ip := range ipSetResp.Accelerator.IpSets {
			for _, publicIP := range ip.IpAddresses {
				results = append(results, ScanResult{
					AccountID:   accountID,
					AccountName: accountName,
					Region:      "global",
					Service:     "globalaccelerator",
					ResourceID:  accName,
					DNSName:     "",
					PublicIP:    publicIP,
					Extra:       "",
					ScanTarget:  publicIP,
					URL:         fmt.Sprintf("https://console.aws.amazon.com/globalaccelerator/home?region=%s#/accelerator/%s", cfg.Region, accName),
				})
			}
		}
	}

	return results
}

