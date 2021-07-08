#!/usr/bin/env bash

set -e

PSQL_CONTAINER="macds-api-postgres-dev"
OUT_CSV="current-90-day-$(date +%Y-%m-%d).csv"

is_container_running() {
  if [ "$(docker container inspect -f '{{.State.Status}}' "$1" 2>/dev/null)" == "running" ]; then
    return 0
  else
    return 1
  fi
}

sql_query() {
  docker exec "${PSQL_CONTAINER}" psql -d postgres -U postgres -c "$1"
}

raw_sql_query() {
  docker exec "${PSQL_CONTAINER}" psql -d postgres -U postgres --quiet --csv -c "$1" | sed -n 2p
}

# Fire up postgresql docker and wait
if ! is_container_running "${PSQL_CONTAINER}"; then
  echo "Postgres is not running. Starting ..."
  bash ./dev_postgres.sh

  TRIES=0
  while ! is_container_running "${PSQL_CONTAINER}"; do
    sleep 1
    echo "Waiting for container to start ..."
    TRIES=$((TRIES + 1))
    if [ $TRIES -gt 15 ]; then
      echo "Postgres wouldn't start"
      exit 2
    fi
  done

  echo "Sleeping for 15 seconds ..."
  sleep 15
fi

# --------- #
# B U I L D #
# --------- #
go build -o macd .

echo -n "Fetching current row count ... "
NUM_HISTORIC=$(raw_sql_query "SELECT COUNT(*) FROM historicals")
echo "${NUM_HISTORIC}"

# --------- #
# F E T C H #
# --------- #

# Fetch Symbols
echo -n "Loading symbol count ... "
NUM_SYMBOLS=$(raw_sql_query "SELECT COUNT(*) FROM symbols")
echo "${NUM_SYMBOLS}"

# Fetch new symbols?
if [ "${NUM_SYMBOLS}" -le "0" ]; then
  echo "  -> No symbols found. Fetching ..."
  ./macd fetch symbols --save
fi

# ---

# Fetch historical data
echo "Fetching historical data ..."
./macd fetch historical --save --gsize 20

echo "Fetching new row count after fetch ... "
raw_sql_query "SELECT COUNT(*) FROM historicals"

# TODO: delete-duplicates.sql is deprecated
# echo "Deleting duplicates ..."
# sql_query "$(cat delete-duplicates.sql)"
# echo "Fetching new row count after dup-deletion ..."
# sql_query "SELECT COUNT(*) FROM historicals"

echo "Writing to csv ..."
docker exec "${PSQL_CONTAINER}" psql -d postgres -U postgres -c \
  "COPY ($(cat 90-day-historical.sql)) TO STDOUT WITH CSV HEADER DELIMITER E'\t'" |
  tr '.' ',' >"${OUT_CSV}"

# overwrite main csv
cat "${OUT_CSV}" >current-90-day.csv
