package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type Transaction struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	FromAccountID uint      `json:"from_account_id"`
	ToAccountID   uint      `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

type BankAccount struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Owner     string  `json:"owner"`
	AccountNo string  `json:"account_no"`
	BankName  string  `json:"bank_name"`
	Balance   float64 `json:"balance"`
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	ToEmail      string
}

func main() {
	log.Println("Starting transaction email service...")

	// Fetch and log accounts for all users in transactions
	err := fetchAndLogUserAccounts()
	if err != nil {
		log.Printf("Warning: Failed to fetch user accounts: %v", err)
	}

	// Get configuration from environment variables
	config := getEmailConfig()

	// Fetch transactions from internal API
	transactions, err := fetchTransactions()
	if err != nil {
		log.Fatalf("Failed to fetch transactions: %v", err)
	}

	log.Printf("Fetched %d transactions from last 24 hours", len(transactions))

	// Generate CSV content
	csvContent, err := generateCSV(transactions)
	if err != nil {
		log.Fatalf("Failed to generate CSV: %v", err)
	}

	// Send email with CSV attachment
	err = sendEmail(config, csvContent, len(transactions))
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}

func getEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsIntOrDefault("SMTP_PORT", 587),
		SMTPUsername: getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", ""),
		ToEmail:      getEnvOrDefault("TO_EMAIL", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func fetchTransactions() ([]Transaction, error) {
	// Get the banking service URL from environment or use default
	baseURL := getEnvOrDefault("CHOREO_INTERNAL_SERVICEURL", "http://localhost:8080")
	url := fmt.Sprintf("%s/transactions", baseURL)

	log.Printf("Fetching transactions from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var transactions []Transaction
	err = json.Unmarshal(body, &transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return transactions, nil
}

func fetchAndLogUserAccounts() error {
	// Get unique user IDs from transactions
	userIDs := make(map[uint]bool)
	userIDs[1] = true // Simulating user IDs for demonstration

	log.Printf("Found %d unique users in transactions", len(userIDs))

	// Fetch accounts for each user
	for userID := range userIDs {
		accounts, err := fetchUserAccounts(userID)
		if err != nil {
			log.Printf("Error fetching accounts for user %d: %v", userID, err)
			continue
		}

		log.Printf("User %d has %d accounts:", userID, len(accounts))
		for _, account := range accounts {
			log.Printf("  - Account ID: %d, Account No: %s, Bank: %s, Owner: %s, Balance: %.2f %s",
				account.ID, account.AccountNo, account.BankName, account.Owner, account.Balance, "USD")
		}
	}

	return nil
}

func fetchUserAccounts(userID uint) ([]BankAccount, error) {
	baseURL := "http://backend-1415777742:8080"
	url := fmt.Sprintf("%s/public/users/%d/accounts", baseURL, userID)

	log.Printf("Fetching accounts for user %d from: %s", userID, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status code: %d", resp.StatusCode)
		return fetchUserAccountsFromPublicAPI(userID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var accounts []BankAccount
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return accounts, nil
}

func fetchUserAccountsFromPublicAPI(userID uint) ([]BankAccount, error) {
	baseURL := "https://7a612bf1-3a19-48ae-8ae3-6a091f5b1370-dev.e1-us-east-azure.choreoapis.dev/banking/backend/v1.0"
	url := fmt.Sprintf("%s/users/%d/accounts", baseURL, userID)

	log.Printf("Fetching accounts for user %d from: %s", userID, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var accounts []BankAccount
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return accounts, nil
}

func generateCSV(transactions []Transaction) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{"ID", "User ID", "From Account ID", "To Account ID", "Amount", "Currency", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write transaction data
	for _, tx := range transactions {
		record := []string{
			strconv.FormatUint(uint64(tx.ID), 10),
			strconv.FormatUint(uint64(tx.UserID), 10),
			strconv.FormatUint(uint64(tx.FromAccountID), 10),
			strconv.FormatUint(uint64(tx.ToAccountID), 10),
			fmt.Sprintf("%.2f", tx.Amount),
			tx.Currency,
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

func sendEmail(config *EmailConfig, csvContent []byte, transactionCount int) error {
	// Validate required configuration
	if config.SMTPUsername == "" || config.SMTPPassword == "" || config.FromEmail == "" || config.ToEmail == "" {
		return fmt.Errorf("missing required email configuration. Please set SMTP_USERNAME, SMTP_PASSWORD, FROM_EMAIL, and TO_EMAIL environment variables")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.FromEmail)
	m.SetHeader("To", config.ToEmail)
	m.SetHeader("Subject", fmt.Sprintf("Daily Transaction Report - %s", time.Now().Format("2006-01-02")))

	// Email body
	body := fmt.Sprintf(`
Dear Team,

Please find attached the transaction report for the last 24 hours.

Summary:
- Total Transactions: %d
- Report Generated: %s
- Time Period: Last 24 hours

Best regards,
Banking System
`, transactionCount, time.Now().Format("2006-01-02 15:04:05"))

	m.SetBody("text/plain", body)

	// Attach CSV file
	fileName := fmt.Sprintf("transactions_%s.csv", time.Now().Format("2006-01-02"))
	m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(csvContent)
		return err
	}))

	// Send email
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	log.Printf("Sending email to: %s", config.ToEmail)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
