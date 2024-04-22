#!/bin/bash
set -e

echo "creating schema"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mydatabase4_2" <<-EOSQL

    CREATE TABLE IF NOT EXISTS public.toy_dynamo (
        id SERIAL PRIMARY KEY,
        key VARCHAR(255) NOT NULL UNIQUE,
        value TEXT NOT NULL,
        expired_at timestamp NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_toy_dynamo_expired_at ON public.toy_dynamo(expired_at);

        
EOSQL

