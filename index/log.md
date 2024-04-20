``` sh
❯ make psql
docker compose exec "postgres" psql -U "ubuntu" "posgre"
psql (15.2)
Type "help" for help.

\d
Did not find any relations.

CREATE TABLE cards (
   mark varchar(255),
   number varchar(255)
);

CREATE TABLE cards (
   mark char(255),
   number char(255)
);

CREATE TABLE cards (
   mark char(4098),
   number char(4099)
);
CREATE INDEX idx_cards ON cards(mark, number);


CREATE TABLE

CREATE INDEX idx_cards ON cards(mark, number);


CREATE INDEX
EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';
                               QUERY PLAN                               
------------------------------------------------------------------------
 Bitmap Heap Scan on cards  (cost=4.18..12.64 rows=4 width=64)
   Recheck Cond: (mark = 'asdf'::text)
   ->  Bitmap Index Scan on idx_cards  (cost=0.00..4.18 rows=4 width=0)
         Index Cond: (mark = 'asdf'::text)
(4 rows)

EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';
                      QUERY PLAN                      
------------------------------------------------------
 Seq Scan on cards  (cost=0.00..1.65 rows=13 width=8)
   Filter: (mark = 'asdf'::text)
(2 rows)

ANALYZE cards;
ANALYZE

EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';
                      QUERY PLAN                      
------------------------------------------------------
 Seq Scan on cards  (cost=0.00..1.65 rows=13 width=8)
   Filter: (mark = 'asdf'::text)
(2 rows)
```

```
SHOW max_connections;
SHOW shared_buffers;
SHOW effective_cache_size;
SHOW maintenance_work_mem;
SHOW checkpoint_completion_target;
SHOW wal_buffers;
SHOW default_statistics_target;
SHOW random_page_cost;
SHOW effective_io_concurrency;
SHOW work_mem;
SHOW huge_pages;
SHOW min_wal_size;
SHOW max_wal_size;

```

変更前の値。

``` sql
SHOW max_connections;
W checkpoint_completion_target;
SHOW wal_buffers;
 max_connections 
-----------------
 100
(1 row)

SHOW shared_buffers;
 shared_buffers 
----------------
 128MB
(1 row)

SHOW effective_cache_size;
 effective_cache_size 
----------------------
 4GB
(1 row)

SHOW maintenance_work_mem;
 maintenance_work_mem 
----------------------
 64MB
(1 row)

SHOW checkpoint_completion_target;
 checkpoint_completion_target 
------------------------------
 0.9
(1 row)

SHOW wal_buffers;
 wal_buffers 
-------------
 4MB
(1 row)

SHOW default_statistics_target;
 default_statistics_target 
---------------------------
 100
(1 row)

SHOW random_page_cost;
 random_page_cost 
------------------
 4
(1 row)

SHOW effective_io_concurrency;
 effective_io_concurrency 
--------------------------
 1
(1 row)

SHOW work_mem;
 work_mem 
----------
 4MB
(1 row)

SHOW huge_pages;
 huge_pages 
------------
 try
(1 row)

SHOW min_wal_size;
 min_wal_size 
--------------
 80MB
(1 row)

SHOW max_wal_size;
 max_wal_size 
--------------
 1GB
(1 row)
```


```
ALTER SYSTEM SET
 min_wal_size = '80MB';


ALTER SYSTEM SET
 shared_buffers = '1MB';
```



``` 
SELECT pg_reload_conf();

ANALYZE cards;

EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';
```


```
/*+ IndexScan(cards idx_cards) */
EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';
```
