#!/usr/bin/env bash

set -e

PSQL_CONTAINER="macds-api-postgres-dev"
OUT_CSV="current-90-day-$(date +%Y-%m-%d).csv"

# ------------- #
# O P T I O N S #
# ------------- #
SKIP_SYMBOL_FETCH=false
SKIP_HISTORIC_FETCH=false
SKIP_DUPE_CHECK=false

while [ "$1" != "" ]; do
  case $1 in
  --skip-symbol)
    SKIP_SYMBOL_FETCH=true
    ;;
  --skip-historic)
    SKIP_HISTORIC_FETCH=true
    ;;
  --skip-dupe)
    SKIP_DUPE_CHECK=true
    ;;
  esac
  shift
done

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

echo -n "[📚] Fetching current row count: "
raw_sql_query "SELECT COUNT(*) FROM historicals"

# --------- #
# F E T C H #
# --------- #

# Fetch Symbols
echo -n "[🏷] Loading symbol count: "
NUM_SYMBOLS=$(raw_sql_query "SELECT COUNT(*) FROM symbols")
echo "${NUM_SYMBOLS}"

# Fetch new symbols?
if [ "${NUM_SYMBOLS}" -le "0" ]; then
  if [ ${SKIP_SYMBOL_FETCH} == false ]; then
    echo "     -> No symbols found. Fetching ..."
    ./macd fetch symbols --save
  else
    echo "     -> Skipped fetching symbols!"
  fi
fi

# ---

# Fetch historical data
if [ ${SKIP_HISTORIC_FETCH} == false ]; then
  echo "[💡] Fetching historical data ..."
  ./macd fetch historical --save --gsize 20
else
  echo "[💡] Skipped fetching historical data!"
fi

# ---

echo -e "[📚] Fetching new row count after fetch ... "
raw_sql_query "SELECT COUNT(*) FROM historicals"

# ------- #
# D U P E #
# ------- #

if [ ${SKIP_DUPE_CHECK} == false ]; then
  echo "[🗑] Deleting Duplicates. (CURRENTLY DEPRECATED. DOES NOT DO ANYTHING!)" # TODO
  echo -n "  [🗑] Before: "
  raw_sql_query "SELECT COUNT(*) FROM historicals"

  # sql_query "$(cat delete-duplicates.sql)"

  echo -n "  [🗑] After: "
  raw_sql_query "SELECT COUNT(*) FROM historicals"
fi

# ----------- #
# E X P O R T #
# ----------- #

echo "[📖] Writing to csv ..."
docker exec "${PSQL_CONTAINER}" psql -d postgres -U postgres -c \
  "COPY ($(cat 90-day-historical.sql)) TO STDOUT WITH CSV HEADER DELIMITER E'\t'" |
  tr '.' ',' >"${OUT_CSV}"

# overwrite main csv
cat "${OUT_CSV}" >current-90-day.csv
