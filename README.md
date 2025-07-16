# AWS Org Public Scanner

`aws-org-scanner` is a command-line tool that discovers publicly accessible resources across AWS accounts, regions, and services. It supports both individual account scanning and AWS Organizations. The tool is fast, concurrent, and provides CSV output for further analysis.

## Overview

The scanner inspects AWS services to identify public-facing resources such as:

- EC2 public IPs
- Load balancers (ELB, ALB, NLB)
- CloudFront distributions
- AppRunner services
- Lambda function URLs
- Route53 public zones
- OpenSearch domains
- RDS instances
- Redshift clusters
- Elastic IPs
- S3 buckets
- AppSync APIs
- Amplify apps
- Global Accelerator

The goal is to produce a list of endpoints (public IPs or URLs) that could be used for further security review or attack surface reduction.

## Why Scan Public Resources?

Exposed cloud infrastructure can lead to data leaks, breaches, or lateral movement. This tool identifies:

| Service        | Why it matters                                             |
|----------------|------------------------------------------------------------|
| EC2            | Instances with public IPs may expose SSH, RDP, or apps     |
| ELB            | Load balancers forward traffic to internal networks        |
| S3             | Public buckets may leak sensitive data                     |
| CloudFront     | Can expose origin resources or leak via misconfigurations  |
| Lambda URLs    | Public serverless endpoints with direct access             |
| AppRunner      | Managed apps with exposed HTTPS endpoints                  |
| RDS            | Public databases can be queried if exposed                 |
| OpenSearch     | Open endpoints can expose logs or data                     |
| Redshift       | Public clusters allow data warehouse access                |
| EIP            | Reserved public IPs may be used or left open               |
| Route53        | Public zones can show exposed DNS records                  |
| Global Accel.  | Globally exposed entry points to regional services         |

## Features

- Scan single AWS account or entire AWS Organization
- Concurrent execution per account, region, and service
- Outputs structured CSV report with public endpoints
- Auto-discovers enabled AWS regions
- Uses `OrganizationAccountAccessRole` or custom IAM roles
- Supports service filtering with `--services`

## Prerequisites

- **Go** 1.20 or higher
- **AWS CLI** configured
- Sufficient IAM permissions to describe resources and (if using `--org`) to assume roles in member accounts

Recommended permissions:
- `sts:AssumeRole`
- `organizations:ListAccounts`
- `ec2:Describe*`, `rds:Describe*`, etc.

## Installation

```bash
git clone https://github.com/your-username/aws-org-scanner.git
cd aws-org-scanner
go build -o scanner ./cmd/scanner
```

## Usage

Scan all accounts in your AWS Organization:

```bash
./scanner --org
```

Scan a single account by ID:

```bash
./scanner --account 123456789012
```

Specify services (comma-separated):

```bash
./scanner --org --services ec2,s3,cloudfront
```

Custom output file:

```bash
./scanner --org --output results.csv
```

Use a custom IAM role name:

```bash
./scanner --org --role CustomAuditRole
```

## Sample Output

Output is a CSV file with the following columns:

| Account ID | Account Name | Region | Service | Resource ID | DNS Name | Public IP | Extra | Scan Target | URL |
|------------|--------------|--------|---------|-------------|----------|-----------|-------|--------------|-----|
| 123456789012 | dev-account | us-east-1 | EC2 | i-0abc123def456 | - | 3.91.203.15 | - | 3.91.203.15 | - |
| 123456789012 | dev-account | global | CloudFront | ABCD1234 | d1234.cloudfront.net | - | - | d1234.cloudfront.net | https://d1234.cloudfront.net |

The `Scan Target` field is either an IP or hostname and can be used for vulnerability scans or traffic monitoring.

## Disclaimer

This tool and associated scripts were created as a personal learning space to explore cloud security concepts and build proof-of-concepts (POCs).  
You are responsible for testing and validating the code before using it in production or enterprise environments.

Please review the code and share feedback or improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
