package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	scan "aws-org-scanner/internal/scan"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const defaultRoleName = "OrganizationAccountAccessRole"

func main() {
	start := time.Now()
	ctx := context.Background()

	orgMode := flag.Bool("org", false, "Scan all accounts in the organization")
	accountID := flag.String("account", "", "Scan a single AWS account (uses current credentials)")
	outputFile := flag.String("output", "scan.csv", "CSV output file name")
	services := flag.String("services", "all", "Comma-separated list of services or 'all'")
	roleName := flag.String("role", defaultRoleName, "IAM Role name to assume in target accounts")

	flag.Parse()

	if !*orgMode && *accountID == "" {
		log.Fatalf("❌ Please specify either --org or --account")
	}

	serviceFilter := make(map[string]bool)
	for _, svc := range strings.Split(*services, ",") {
		serviceFilter[strings.ToLower(strings.TrimSpace(svc))] = true
	}

	rootCfg, err := scan.LoadAWSConfig(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}

	rootAccountID := getCallerAccountID(ctx, rootCfg)

	var accounts []scan.Account
	if *orgMode {
		accounts, err = scan.GetOrgAccounts(ctx, rootCfg)
		if err != nil {
			log.Fatalf("❌ Failed to retrieve org accounts: %v", err)
		}
	} else {
		accounts = []scan.Account{{ID: *accountID, Name: "SingleAccount"}}
	}

	fmt.Printf("🔍 Scanning %d account(s)...\n", len(accounts))

	var allResults []scan.ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, account := range accounts {
		wg.Add(1)
		go func(account scan.Account) {
			defer wg.Done()

			var cfg aws.Config
			if *orgMode && account.ID != rootAccountID {
				var err error
				cfg, err = scan.AssumeRole(ctx, rootCfg, account.ID, *roleName)
				if err != nil {
					log.Printf("⚠️  Failed to assume role for account %s: %v\n", account.ID, err)
					return
				}
			} else {
				cfg = rootCfg
			}

			regions := scan.GetEnabledRegions(ctx, cfg)

			var regionWG sync.WaitGroup
			for _, region := range regions {
				regionWG.Add(1)
				go func(region string) {
					defer regionWG.Done()
					results := scan.ScanRegion(ctx, cfg, account.ID, account.Name, region, serviceFilter)
					mu.Lock()
					allResults = append(allResults, results...)
					mu.Unlock()
				}(region)
			}
			regionWG.Wait()
		}(account)
	}

	wg.Wait()

	err = writeCSV(*outputFile, allResults)
	if err != nil {
		log.Fatalf("❌ Failed to write output: %v", err)
	}

	fmt.Printf("✅ Scan complete. Found %d public resources. Output saved to %s\n", len(allResults), *outputFile)
	fmt.Printf("⏱️  Total time: %s\n", time.Since(start).Round(time.Second))
}

func getCallerAccountID(ctx context.Context, cfg aws.Config) string {
	stsClient := sts.NewFromConfig(cfg)
	resp, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("❌ Failed to get caller identity: %v", err)
	}
	return aws.ToString(resp.Account)
}

func writeCSV(filename string, results []scan.ScanResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Account ID",
		"Account Name",
		"Region",
		"Service",
		"Resource ID",
		"DNS Name",
		"Public IP",
		"Extra",
		"Scan Target",
		"URL",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, r := range results {
		row := []string{
			r.AccountID,
			r.AccountName,
			r.Region,
			r.Service,
			r.ResourceID,
			r.DNSName,
			r.PublicIP,
			r.Extra,
			r.ScanTarget,
			r.URL,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

