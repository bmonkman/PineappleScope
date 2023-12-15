insert into firings VALUES (NULL, "2018-03-07 07:00:29", "2018-03-07 10:55:29", 6.66, 1, "Test Data", "Testing", 0, 100, 0, 0);

WITH RECURSIVE cte AS (
    SELECT 1 AS n
    UNION ALL
    SELECT n + 1 FROM cte WHERE n < 100
)
INSERT INTO temperature_readings
SELECT NULL, datetime('2018-03-07 07:00:29', '+' || (n - 1) || ' minutes'), (SELECT id FROM firings WHERE name='Test Data'), ROUND(10.00 + (n - 1) * 0.25 - 12.5 * SIN((n - 1) * 0.1), 2), 6.66
FROM cte;
