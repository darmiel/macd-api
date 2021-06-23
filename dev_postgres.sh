#!/bin/bash

docker run -d --rm \
    --name macds-api-postgres-dev \
    -e POSTGRES_PASSWORD=123456 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "$(pwd)/data/postgres:/var/lib/postgresql/data" \
    -p 45432:5432 \
    postgres