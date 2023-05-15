--- db size ---
SELECT
    pg_size_pretty (
        pg_database_size ('nogi-official')
    );
--  pg_size_pretty 
-- ----------------
--  2396 MB
-- (1 row)


--- table size ---
SELECT
    relname AS "relation",
    pg_size_pretty (
        pg_total_relation_size (C .oid)
    ) AS "total_size"
FROM
    pg_class C
LEFT JOIN pg_namespace N ON (N.oid = C .relnamespace)
WHERE
    nspname NOT IN (
        'pg_catalog',
        'information_schema'
    )
AND C .relkind <> 'i'
AND nspname !~ '^pg_toast'
ORDER BY
    pg_total_relation_size (C .oid) DESC
LIMIT 5;

--           relation           | total_size 
-- -----------------------------+------------
--  comments                    | 2292 MB
--  blog_images                 | 71 MB
--  blogs                       | 24 MB
--  comment_counts_rank_by_blog | 1192 kB
--  members                     | 72 kB
-- (5 rows)
