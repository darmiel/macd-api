WITH rank AS (
    SELECT h.symbol,
           h.day_date,
           ROW_NUMBER() OVER (PARTITION BY h.symbol,
			to_char(h.day_date,
			'YYYYMMDD')
			ORDER BY
				h.day_date DESC) AS RowRank -- DESC: 18:29 -> 2; 18:32 -> 1
    FROM historicals h
)
DELETE FROM historicals t
WHERE t.id IN (
    SELECT
    h.id FROM historicals h
    INNER JOIN rank ON rank.symbol = h.symbol
  AND rank.day_date = h.day_date
    WHERE
    rank.RowRank != 1)