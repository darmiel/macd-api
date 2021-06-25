WITH a AS (
    SELECT symbol
    FROM historicals
    GROUP BY symbol
    HAVING COUNT(symbol) >= 90
)

SELECT
       *
FROM (
         SELECT
                ROW_NUMBER() OVER (PARTITION BY h.symbol ORDER BY h.timestamp DESC) AS r,
                h.symbol,
                to_timestamp(h.timestamp) AS date,
                h.high,
                h.low,
                h.open,
                h.close,
                h.volume
         FROM
             historicals h
             INNER JOIN a
         ON a.symbol = h.symbol
     ) x
WHERE x.r <= 90