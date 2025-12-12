#!/bin/bash

# ============================================================================
# Health Check Script
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Checking services health...${NC}"
echo ""

# Check PostgreSQL
if docker exec alerting-postgres pg_isready -U postgres -d alerting_db > /dev/null 2>&1; then
    echo -e "PostgreSQL: ${GREEN}✅ Healthy${NC}"
else
    echo -e "PostgreSQL: ${RED}❌ Not healthy${NC}"
fi

# Check Redis
if docker exec alerting-redis redis-cli ping > /dev/null 2>&1; then
    echo -e "Redis:      ${GREEN}✅ Healthy${NC}"
else
    echo -e "Redis:      ${RED}❌ Not healthy${NC}"
fi

# Check API
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "API:        ${GREEN}✅ Healthy${NC}"
else
    echo -e "API:        ${YELLOW}⚠️  Not responding (may not be running)${NC}"
fi

echo ""
