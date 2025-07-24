# Transaction Email Service

This service fetches transactions from the last 24 hours using the internal banking API and sends them as a CSV attachment via email.

## Configuration

The service uses environment variables for configuration:

### Required Variables
- `SMTP_USERNAME`: SMTP username for email authentication
- `SMTP_PASSWORD`: SMTP password for email authentication  
- `FROM_EMAIL`: Email address to send from
- `TO_EMAIL`: Email address to send to

### Optional Variables
- `SMTP_HOST`: SMTP server host (default: smtp.gmail.com)
- `SMTP_PORT`: SMTP server port (default: 587)
- `BANKING_SERVICE_URL`: Base URL of the banking service (default: http://localhost:8080)

## Usage

### Using Gmail SMTP

1. Set up environment variables:
```bash
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"  # Use App Password for Gmail
export FROM_EMAIL="your-email@gmail.com"
export TO_EMAIL="recipient@example.com"
export BANKING_SERVICE_URL="http://localhost:8080"
```

2. Run the service:
```bash
go run main.go
```

### Building and Running

```bash
# Build the binary
go build -o transaction-emailer main.go

# Run the binary
./transaction-emailer
```

## How it works

1. The service calls the `/internal/transactions` endpoint of the banking service
2. Fetches all transactions from the last 24 hours
3. Converts the transaction data to CSV format
4. Sends an email with the CSV file as an attachment
5. Includes a summary in the email body

## CSV Format

The CSV file contains the following columns:
- ID: Transaction ID
- User ID: User who made the transaction
- From Account ID: Source account ID
- To Account ID: Destination account ID
- Amount: Transaction amount
- Currency: Transaction currency
- Created At: Transaction timestamp

## Scheduling

This service is designed to be run as a scheduled task (e.g., daily cron job). You can schedule it using:

### Cron (Linux/macOS)
```bash
# Run daily at 9 AM
0 9 * * * /path/to/transaction-emailer
```

### Task Scheduler (Windows)
Create a scheduled task to run the executable daily.

## Error Handling

The service will log errors and exit with a non-zero status code if:
- Unable to fetch transactions from the API
- Unable to generate CSV content
- Unable to send email
- Missing required configuration

## Security Notes

- Use app passwords instead of regular passwords for Gmail
- Store sensitive credentials securely
- Consider using environment files or secret management systems in production
