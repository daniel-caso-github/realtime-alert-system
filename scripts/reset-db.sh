#!/bin/bash

# ============================================================================
# Reset Database Script
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}⚠️  This will delete all data in the database!${NC}"
read -p "Are you sure? (y/N) " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Aborted.${NC}"
    exit 1
fi

echo -e "${YELLOW}Stopping containers...${NC}"
docker-compose down

echo -e "${YELLOW}Removing PostgreSQL volume...${NC}"
docker volume rm alerting-postgres-data 2>/dev/null || true

echo -e "${YELLOW}Starting PostgreSQL...${NC}"
docker-compose up -d postgres

echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
sleep 5

echo -e "${GREEN}✅ Database reset complete!${NC}"
echo -e "${GREEN}Run 'docker-compose up -d' to start all services.${NC}"
