#!/bin/bash

DOCKER_CONTAINER="macds-api-postgres-dev"

docker exec "${DOCKER_CONTAINER}" \
    psql -d postgres -U postgres -c \
    "COPY ($(cat 90-day-historical.sql)) TO STDOUT WITH CSV HEADER DELIMITER E'\t'" | tr '.' ','
