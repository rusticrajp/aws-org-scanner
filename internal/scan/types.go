package scan

// Account holds basic account ID and name.
type Account struct {
	ID   string
	Name string
}

// ScanResult represents a public-facing resource found in the scan.
type ScanResult struct {
	AccountID   string
	AccountName string
	Region      string
	Service     string
	ResourceID  string
	DNSName     string
	Extra       string
	PublicIP    string // ✅ Added for IP-addressed resources (EC2, RDS, etc.)
	ScanTarget  string // IP, DNS, or URL suitable for scanning
	URL         string // Optional full URL (e.g. for AppRunner, Amplify, etc.)
}

