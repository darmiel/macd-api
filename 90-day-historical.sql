WITH a AS (
    SELECT symbol
    FROM historics
    GROUP BY symbol
    HAVING COUNT(symbol) >= 90
)
SELECT x.* -- , s.use
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
         FROM historics h
                  INNER JOIN a
                             ON a.symbol = h.symbol) x
         INNER JOIN symbols s ON s.symbol = x.symbol
WHERE x.r <= 90
  AND s.use = TRUE