INSERT INTO statistics ("tgUserId", "exercise", "count", "params", "createdAt", "statusId")
SELECT
    -- Generate random tgUserId (simulate ~10000 different users)
    ('user_' || (random() * 1000)::int)::varchar(255),

    -- Generate random exercises from a realistic set
    (ARRAY[
         'pullUp',
     'muscleUp',
     'pushUp',
     'dip',
     'abs',
     'squat',
     'lunge',
     'burpee',
     'skippingRope'
         ])[floor(random() * 8 + 1)],

    -- Generate random count (1-50 reps)
    (random() * 50 + 1)::float8,

    -- Generate some random JSONB params
    '{}',

    -- Generate random timestamps over the last year
    now() - (random() * interval '365 days'),

    -- Generate statusId (98% active, 2% inactive for realistic distribution)
    CASE WHEN random() < 0.98 THEN 1 ELSE 0 END
FROM generate_series(1, 1000000);

ANALYZE statistics;
SELECT
    pg_size_pretty(pg_total_relation_size('statistics')) as table_size,
    count(*) as row_count
FROM statistics;

explain analyze
SELECT t."tgUserId", t."exercise", sum(t.count) as "sumCount", count(t.count) as "sets"
FROM statistics t
WHERE true
  AND t."statusId" = 1
  AND t."tgUserId" = 'user_745'
  AND t."exercise" in ('pullUp', 'pushUp')
  AND (false
    OR (
           t."createdAt" >= '2025-06-01 00:00:00+00:00:00'
               AND t."createdAt" < '2025-06-10 19:29:42.227219+00:00:00'
           )

    OR (
           t."createdAt" >= '2025-05-01 00:00:00+00:00:00'
               AND t."createdAt" < '2025-05-22 20:29:43.227219+00:00:00'
           )
    )
GROUP BY 1, 2
ORDER BY 3 DESC;
