WITH a AS (
    SELECT symbol
    FROM historicals
    GROUP BY symbol
    HAVING COUNT(symbol) >= 90
)
SELECT *
FROM (
         SELECT ROW_NUMBER() OVER (PARTITION BY h.symbol ORDER BY h.day_date DESC) AS r,
                h.symbol,
                h.day_date,
                h.orig_date,
                h.high,
                h.low,
                h.open,
                h.close,
                h.volume
         FROM historicals h
                  INNER JOIN a
                             ON a.symbol = h.symbol) x
WHERE x.r <= 90