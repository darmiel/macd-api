WITH rank AS (
    SELECT h.symbol,
           h.date,
           ROW_NUMBER() OVER (PARTITION BY h.symbol,
			to_char(h.date,
			'YYYYMMDD')
			ORDER BY
				h.date DESC) AS RowRank -- DESC: 18:29 -> 2; 18:32 -> 1
    FROM historicals h
)
DELETE FROM historicals t
WHERE t.id IN (
    SELECT
    h.id FROM historicals h
    INNER JOIN rank ON rank.symbol = h.symbol
  AND rank.date = h.date
    WHERE
    rank.RowRank != 1)