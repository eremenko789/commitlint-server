#!/bin/bash
#
# Health check script for commitlint webhook server
# Can be used with monitoring systems like Nagios, Zabbix, etc.
#

set -e

# Configuration
WEBHOOK_URL="${WEBHOOK_URL:-http://localhost:3000}"
TIMEOUT="${TIMEOUT:-5}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo -e "${RED}ERROR: curl is not installed${NC}"
    exit 2
fi

# Perform health check
echo -n "Checking webhook server at ${WEBHOOK_URL}/health... "

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    --max-time ${TIMEOUT} \
    "${WEBHOOK_URL}/health" 2>/dev/null || echo "000")

case $HTTP_CODE in
    200)
        echo -e "${GREEN}OK${NC}"
        echo "Server is healthy (HTTP ${HTTP_CODE})"
        exit 0
        ;;
    000)
        echo -e "${RED}FAILED${NC}"
        echo "Connection failed or timeout"
        exit 2
        ;;
    *)
        echo -e "${YELLOW}WARNING${NC}"
        echo "Unexpected response (HTTP ${HTTP_CODE})"
        exit 1
        ;;
esac