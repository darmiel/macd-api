DELETE
FROM historicals h
WHERE to_char(h.date, 'YYYY-MM-dd') = to_char(CURRENT_DATE, 'YYYY-MM-dd')