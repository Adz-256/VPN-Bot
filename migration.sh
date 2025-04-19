#!/bin/bash
source .env

export MIGRATION_DSN="${DSN}"
export MIGRATION_DIR="migrations"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v