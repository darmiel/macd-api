WITH rank AS (
    SELECT h.symbol,
           h.day_date,
           ROW_NUMBER() OVER (PARTITION BY h.symbol,
               TO_CHAR(h.day_date,
                       'YYYYMMDD')
               ORDER BY
                   h.day_date DESC) AS RowRank -- DESC: 18:29 -> 2; 18:32 -> 1
    FROM historics h
)

DELETE
FROM historics HIST
    USING rank RA
WHERE HIST.symbol = RA.symbol
  AND HIST.day_date = RA.day_date
  AND RA.RowRank != 1