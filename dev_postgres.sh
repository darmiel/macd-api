#!/bin/bash

docker run -d --rm \
    --name macds-api-postgres-dev \
    -e POSTGRES_PASSWORD=123456 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v "$(pwd)/data/postgres:/var/lib/postgresql/data" \
    -p 45432:5432 \
    postgres

# psql -d postgres -U postgres -c "COPY (${QUERY}) TO STDOUT WITH CSV HEADER DELIMITER E'\t'" | tr '.' ',' > /var/lib/postgresql/data/pgdata/output.csv
# docker exec macds-api-postgres-dev psql -d postgres -U postgres -c "COPY (SELECT * FROM historicals LIMIT 10) TO STDOUT WITH CSV HEADER DELIMITER E'\t'" | tr '.' ','
# docker build -f docker/update-90-days/Dockerfile .