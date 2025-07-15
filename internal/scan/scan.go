package scan

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// ScanRegion runs all service scanners for a given region
func ScanRegion(ctx context.Context, cfg aws.Config, accountID, accountName, region string, serviceFilter map[string]bool) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	cfg.Region = region

	scanFuncs := map[string]func(context.Context, aws.Config, string, string, string) []ScanResult{
		"ec2":               ScanEC2,
		"eip":               ScanEIP,
		"rds":               ScanRDS,
		"elb":               ScanELB,
		"apigateway":        ScanAPIGateway,
		"redshift":          ScanRedshift,
		"opensearch":        ScanOpenSearch,
		"lightsail":         ScanLightsail,
		"elasticbeanstalk":  ScanElasticBeanstalk,
		"cloudfront":        ScanCloudFront,
		"s3":                ScanS3,
		"globalaccelerator": ScanGlobalAccelerator,
		"route53":           ScanRoute53,
		"amplify":           ScanAmplify,
		"appsync":           ScanAppSync,
		"apprunner":         ScanAppRunner,
		"workspaces":        ScanWorkSpaces,
		"lambda":            ScanLambda, // ✅ New Lambda Function URL scanner
	}

	for name, fn := range scanFuncs {
		nameCopy := name
		fnCopy := fn

		if len(serviceFilter) > 0 && !serviceFilter["all"] && !serviceFilter[strings.ToLower(nameCopy)] {
			continue
		}

		if isGlobalService(nameCopy) && region != "us-east-1" {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			svcResults := fnCopy(ctx, cfg, accountID, accountName, region)
			mu.Lock()
			results = append(results, svcResults...)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return results
}

// isGlobalService returns true if a service is account-global (runs only once in us-east-1)
func isGlobalService(name string) bool {
	switch name {
	case "cloudfront", "s3", "globalaccelerator", "route53":
		return true
	default:
		return false
	}
}

