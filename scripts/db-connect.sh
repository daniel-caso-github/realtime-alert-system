#!/bin/bash

# ============================================================================
# Database Connection Script
# ============================================================================

docker exec -it alerting-postgres psql -U postgres -d alerting_db
