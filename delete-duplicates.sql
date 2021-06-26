DELETE
FROM historicals
WHERE id IN (
    SELECT id
    FROM historicals c
             INNER JOIN (
        SELECT h.symbol,
               h.date,
               ROW_NUMBER() OVER (PARTITION BY h.symbol ORDER BY h.date DESC) AS RowRank
        FROM historicals h
                 INNER JOIN (
            SELECT symbol,
                   to_char(date, 'YYYYMMDD') dt,
                   COUNT(*) AS               CountOf
            FROM historicals
            GROUP BY symbol, dt
            HAVING COUNT(*) > 1) dt ON h.symbol = dt.symbol) e ON c.symbol = e.symbol
        AND c.date = e.date
    WHERE e.RowRank != 1)