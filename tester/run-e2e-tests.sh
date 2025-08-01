#!/bin/bash

# Banking API E2E Test Runner
# This script helps you run the Postman collection locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
BASE_URL=""
CHOREO_API_KEY=""
ENVIRONMENT_FILE="environment.json"
VERBOSE=false

# Function to display usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -u, --base-url URL        Base URL for the API (required)"
    echo "  -k, --api-key KEY         Choreo API key (required)"
    echo "  -e, --env-file FILE       Environment file path (default: environment.json)"
    echo "  -v, --verbose             Enable verbose output"
    echo "  -h, --help                Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -u https://api.example.com -k your-api-key"
    echo "  $0 --base-url https://api.example.com --api-key your-key --verbose"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--base-url)
            BASE_URL="$2"
            shift 2
            ;;
        -k|--api-key)
            CHOREO_API_KEY="$2"
            shift 2
            ;;
        -e|--env-file)
            ENVIRONMENT_FILE="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Validate required parameters
if [[ -z "$BASE_URL" ]]; then
    echo -e "${RED}Error: Base URL is required${NC}"
    usage
    exit 1
fi

if [[ -z "$CHOREO_API_KEY" ]]; then
    echo -e "${RED}Error: Choreo API key is required${NC}"
    usage
    exit 1
fi

# Check if newman is installed
if ! command -v newman &> /dev/null; then
    echo -e "${RED}Error: Newman is not installed${NC}"
    echo "Please install it with: npm install -g newman newman-reporter-htmlextra"
    exit 1
fi

# Create environment file
echo -e "${BLUE}Creating environment file: $ENVIRONMENT_FILE${NC}"
cat > "$ENVIRONMENT_FILE" << EOF
{
  "id": "banking-api-env",
  "name": "Banking API Environment",
  "values": [
    {
      "key": "baseUrl",
      "value": "$BASE_URL",
      "enabled": true
    },
    {
      "key": "choreoApiKey",
      "value": "$CHOREO_API_KEY",
      "enabled": true
    }
  ]
}
EOF

echo -e "${GREEN}Environment file created successfully${NC}"

# Prepare newman command
NEWMAN_CMD="newman run postman-collection.json --environment $ENVIRONMENT_FILE --reporters cli,htmlextra --reporter-htmlextra-export test-results-$(date +%Y%m%d-%H%M%S).html --timeout-request 30000 --color on"

if [[ "$VERBOSE" == true ]]; then
    NEWMAN_CMD="$NEWMAN_CMD --verbose"
fi

# Display configuration
echo -e "${YELLOW}Configuration:${NC}"
echo "  Base URL: $BASE_URL"
echo "  API Key: [REDACTED]"
echo "  Environment File: $ENVIRONMENT_FILE"
echo "  Verbose: $VERBOSE"
echo ""

# Run tests
echo -e "${BLUE}Running Banking API E2E Tests...${NC}"
echo ""

if eval $NEWMAN_CMD; then
    echo ""
    echo -e "${GREEN}✅ All tests passed successfully!${NC}"
    echo -e "${BLUE}📊 Check the generated HTML report for detailed results${NC}"
else
    echo ""
    echo -e "${RED}❌ Some tests failed${NC}"
    echo -e "${YELLOW}📊 Check the generated HTML report for detailed error information${NC}"
    exit 1
fi

# Cleanup
if [[ "$ENVIRONMENT_FILE" == "environment.json" ]]; then
    echo -e "${YELLOW}Cleaning up environment file...${NC}"
    rm -f "$ENVIRONMENT_FILE"
fi

echo -e "${GREEN}Test run completed!${NC}"
