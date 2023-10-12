``` sql
CREATE TABLE cards (
   mark text,
   number text
);

CREATE INDEX idx_cards ON cards(mark, number);


echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('abcd', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('bdef', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('asdf', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('ggrk', '\${0}');\""

echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('abcd', '\${0}\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('bdef', '\${0}\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('asdf', '\${0}\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards VALUES ('ggrk', '\${0}\${0}');\""


ANALYZE cards;


EXPLAIN SELECT * FROM cards WHERE mark = 'asdf';


```

``` sql
CREATE TABLE cards_num (
   mark int,
   number text
);

CREATE INDEX idx_cards_num ON cards_num(mark, number);
ANALYZE cards_num;


echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards_num VALUES (1, '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards_num VALUES (2, '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards_num VALUES (3, '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards_num VALUES (4, '\${0}');\""

```



``` sql
CREATE TABLE cards_seq (
   mark text,
   number text
);

CREATE INDEX idx_cards_seq ON cards_seq(mark, number);


echo {1..13} | xargs -n 1 bash -c "psql --host=localhost --username=ubuntu --dbname=posgre -c \"INSERT INTO cards_seq (mark, number) VALUES ('💚', '\${0}'), ('🍀', '\${0}'), ('🔷', '\${0}'), ('♠️', '\${0}');\""

ANALYZE cards_seq;


EXPLAIN SELECT * FROM cards WHERE mark = '🔷';


```


現状確認

``` sh
posgre=# SHOW shared_buffers;
 shared_buffers 
----------------
 128MB
(1 row)


ALTER SYSTEM SET shared_buffers = '1MB';

ALTER SYSTEM SET random_page_cost = '4';


systemctl restart postgresql



ANALYZE cards;



SELECT pg_reload_conf();



posgre=# show effective_cache_size;
 effective_cache_size 
----------------------
 4GB
(1 row)


ALTER SYSTEM SET effective_cache_size = '1MB';
SELECT pg_reload_conf();
ANALYZE cards;



posgre=# SHOW default_statistics_target;
 default_statistics_target 
---------------------------
 100
(1 row)



find / -name postgresql.conf 2> /dev/null
/var/lib/postgresql/data/postgresql.conf




ALTER SYSTEM SET
 max_connections = '200';
ALTER SYSTEM SET
 shared_buffers = '256MB';
ALTER SYSTEM SET
 effective_cache_size = '768MB';
ALTER SYSTEM SET
 maintenance_work_mem = '64MB';
ALTER SYSTEM SET
 checkpoint_completion_target = '0.9';
ALTER SYSTEM SET
 wal_buffers = '7864kB';
ALTER SYSTEM SET default_statistics_target = '33';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET
 random_page_cost = '1.1';
ALTER SYSTEM SET
 effective_io_concurrency = '200';
ALTER SYSTEM SET
 work_mem = '655kB';
ALTER SYSTEM SET
 huge_pages = 'off';
ALTER SYSTEM SET
 min_wal_size = '1GB';
ALTER SYSTEM SET
 max_wal_size = '4GB';


SHOW wal_buffers;

SELECT pg_reload_conf();
ANALYZE cards;


posgre=# SHOW effective_cache_size;
 effective_cache_size 
----------------------
 1MB
(1 row)

ALTER SYSTEM SET effective_cache_size = '4GB';

random_page_cost = 1.1
ALTER SYSTEM SET random_page_cost = '1.1';

ALTER SYSTEM SET shared_buffers = '8GB';

SHOW shared_buffers;
SHOW effective_cache_size;
SHOW random_page_cost;
SHOW default_statistics_target;




ALTER SYSTEM SET default_statistics_target = '33';
ALTER SYSTEM SET default_statistics_target = '100';


SELECT pg_reload_conf(); ANALYZE cards; EXPLAIN SELECT * FROM cards WHERE mark = '🍀';


ALTER SYSTEM SET shared_buffers = '4GB';
ALTER SYSTEM SET effective_cache_size = '4GB';


SET random_page_cost = 1.1;
SET random_page_cost = 0.8;


EXPLAIN SELECT * FROM cards WHERE mark = '🔷';





posgre=# SELECT current_setting('block_size');
 current_setting 
-----------------
 8192
(1 row)

posgre=# SHOW block_size;
 block_size 
------------
 8192
(1 row)



ALTER SYSTEM SET block_size = '16';

```




### 統計情報で変わった例

``` sql
posgre=# EXPLAIN SELECT * FROM cards WHERE mark = '🍀';
                               QUERY PLAN                               
------------------------------------------------------------------------
 Bitmap Heap Scan on cards  (cost=1.28..5.48 rows=4 width=64)
   Recheck Cond: (mark = '🍀'::text)
   ->  Bitmap Index Scan on idx_cards  (cost=0.00..1.28 rows=4 width=0)
         Index Cond: (mark = '🍀'::text)
(4 rows)

posgre=# 
posgre=# ANALYZE cards;
ANALYZE
posgre=# EXPLAIN SELECT * FROM cards WHERE mark = '🍀';
                      QUERY PLAN                      
------------------------------------------------------
 Seq Scan on cards  (cost=0.00..1.65 rows=13 width=7)
   Filter: (mark = '🍀'::text)
(2 rows)



posgre=# SHOW shared_buffers;
 shared_buffers 
----------------
 8GB
(1 row)

posgre=# SHOW effective_cache_size;
 effective_cache_size 
----------------------
 4GB
(1 row)

posgre=# SHOW random_page_cost;
 random_page_cost 
------------------
 1.1
(1 row)

posgre=# SHOW default_statistics_target;
 default_statistics_target 
---------------------------
 100
(1 row)
```
