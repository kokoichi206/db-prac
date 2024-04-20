``` sql
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'your_table_name'
AND indexname LIKE 'your_table_name_pkey';


SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'comments_pkey';


CREATE TABLE kawaiis (
    generation int,
    name char(46)
);

CREATE TABLE kawaiis (
    generation int,
    first char(46),
    last char(46)
);

```

``` sh
psql --host=localhost --username=ubuntu --dbname=nogi-official -c "INSERT INTO kawaiis VALUES (1, 'pi')";

echo {1..5}{a..z} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO kawaiis VALUES ('\${0:0:1}', '\${0:1:1}');\""


CREATE INDEX idx_kawaiis_generation ON kawaiis(generation, name);

# 9 * 26 * 26 = 6084
echo {1..9}{a..z}{a..z} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=nogi-official -c \"INSERT INTO kawaiis VALUES ('\${0:0:1}', '\${0:1:1}', '\${0:2:1}');\""

echo {1..5}{a..z} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=nogi-official -c \"DELETE FROM kawaiis;\""
```

index をはる

``` sql
nogi-official=# \d kawaiis 
                   Table "public.kawaiis"
   Column   |     Type      | Collation | Nullable | Default 
------------+---------------+-----------+----------+---------
 generation | integer       |           |          | 
 name       | character(46) |           |          | 


nogi-official=# \d kawaiis 
                   Table "public.kawaiis"
   Column   |     Type      | Collation | Nullable | Default 
------------+---------------+-----------+----------+---------
 generation | integer       |           |          | 
 name       | character(46) |           |          | 
Indexes:
    "idx_kawaiis_generation" btree (generation, name)
```

explain をみる

``` sql
ANALYZE kawaiis;

-- 両方でインデックス指定
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5 AND name = 'i';
                                         QUERY PLAN                                         
--------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.27..4.29 rows=1 width=51)
   Index Cond: ((generation = 5) AND (name = 'i'::bpchar))
(2 rows)

-- 1つ目でインデックス指定, なぜか Seq Scan になる。
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5;
                       QUERY PLAN                        
---------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..4.62 rows=26 width=51)
   Filter: (generation = 5)
(2 rows)

-- 2つ目でインデックス指定
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE name = 'i';
                       QUERY PLAN                       
--------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..7.90 rows=9 width=51)
   Filter: (name = 'i'::bpchar)
(2 rows)
```

統計情報を増やす。

``` sql
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5 AND name = 'i';
                                         QUERY PLAN                                         
--------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.28..4.30 rows=1 width=51)
   Index Cond: ((generation = 5) AND (name = 'i'::bpchar))
(2 rows)

nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5;
                                          QUERY PLAN                                           
-----------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.28..60.11 rows=676 width=51)
   Index Cond: (generation = 5)
(2 rows)

nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE name = 'ni';
                        QUERY PLAN                        
----------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..142.05 rows=9 width=51)
   Filter: (name = 'ni'::bpchar)
(2 rows)
```



テーブルを変える

``` sql
CREATE TABLE kawaiis (
    generation int,
    first char(46),
    last char(46)
);
```

``` sh
CREATE INDEX idx_kawaiis_generation ON kawaiis(generation, first, last);

# 9 * 26 * 26 = 6084
echo {1..9}{a..z}{a..z} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=nogi-official -c \"INSERT INTO kawaiis VALUES ('\${0:0:1}', '\${0:1:1}', '\${0:2:1}');\""
```

``` sql
-- 1 and 2 and 3
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5 AND first = 'a' AND last = 'o';
                                         QUERY PLAN                                         
--------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.41..4.43 rows=1 width=98)
   Index Cond: ((generation = 5) AND (first = 'a'::bpchar) AND (last = 'o'::bpchar))
(2 rows)

-- 1 and 2
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5 AND first = 'a';
                                         QUERY PLAN                                          
---------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.41..4.93 rows=26 width=98)
   Index Cond: ((generation = 5) AND (first = 'a'::bpchar))
(2 rows)

-- 1 and 3
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5 AND last = 'o';
                                          QUERY PLAN                                          
----------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.41..55.43 rows=26 width=98)
   Index Cond: ((generation = 5) AND (last = 'o'::bpchar))
(2 rows)     

-- 2 and 3
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE first ='a' AND last = 'o';
                         QUERY PLAN                         
------------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..191.26 rows=9 width=98)
   Filter: ((first = 'a'::bpchar) AND (last = 'o'::bpchar))
(2 rows)

-- 1
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE generation = 5;
                                          QUERY PLAN                                           
-----------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.41..60.24 rows=676 width=98)
   Index Cond: (generation = 5)
(2 rows)

-- 2
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE first ='a';
                         QUERY PLAN                         
------------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..176.05 rows=234 width=98)
   Filter: (first = 'a'::bpchar)
(2 rows)

-- 3
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE last ='o';
                         QUERY PLAN                         
------------------------------------------------------------
 Seq Scan on kawaiis  (cost=0.00..176.05 rows=234 width=98)
   Filter: (last = 'o'::bpchar)
(2 rows)


-- 2 and 1
nogi-official=# EXPLAIN SELECT * FROM kawaiis WHERE first = 'a' AND generation = 5;
                                         QUERY PLAN                                          
---------------------------------------------------------------------------------------------
 Index Only Scan using idx_kawaiis_generation on kawaiis  (cost=0.41..4.93 rows=26 width=98)
   Index Cond: ((generation = 5) AND (first = 'a'::bpchar))
(2 rows)
```



``` sql
EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM kawaiis WHERE generation = 5 AND first = 'a';


```

