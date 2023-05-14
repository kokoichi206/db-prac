## VIEW, MATERIALIZED VIEW

一般的な VIEW と MATERIALIZED VIEW の主な違いは、MATERIALIZED VIEW が物理的にデータを格納する点にあります。
これは一般的なVIEWがクエリを保存するだけでデータを保存しないのとは対照的です。

MATERIALIZED VIEWはデータベース内に実際のデータを格納します。つまり、MATERIALIZED VIEWは指定されたクエリを実行し、その結果をデータベース内に格納します。そのため、MATERIALIZED VIEWは一種のキャッシュとして機能します。MATERIALIZED VIEWへのクエリは、データベースから直接データを取得するのではなく、MATERIALIZED VIEWに格納されているデータを使用するため、パフォーマンスの向上をもたらすことがあります。

ただし、MATERIALIZED VIEWのデータは静的であるため、基になるテーブルのデータが変更されると、それに応じてMATERIALIZED VIEWも手動または定期的に更新（リフレッシュ）する必要があります。そうしないと、MATERIALIZED VIEWが古いデータを表示し続け、その結果と基本テーブルのデータとの間に不一致が生じます。

## how to create materialized viwe

PostgreSQL で MATERIALIZED VIEW を作成するには、以下のようなSQLコマンドを使用します

``` sql
CREATE MATERIALIZED VIEW view_name AS
SELECT column1, column2, ...
FROM table_name
WHERE condition;
```

### やったこと

``` sql
-- materialized view
CREATE MATERIALIZED VIEW comment_counts_rank_by_blog AS
SELECT
    c.blog_code,
    COUNT(c) AS count,
    SUBSTRING(b.title FROM 1 FOR 30)  AS title,
    to_timestamp(b.timestamp::bigint / 1000000000) AS date,
    m.name AS name,
    b.link AS link
FROM comments AS c
INNER JOIN blogs AS b ON c.blog_code = b.code
INNER JOIN members AS m ON b.member_code = m.code
GROUP BY c.blog_code, b.timestamp, m.name, b.link, b.title
ORDER BY count DESC;


nogi-official=# \d
                         List of relations
 Schema |            Name             |       Type        | Owner  
--------+-----------------------------+-------------------+--------
 public | blog_images                 | table             | ubuntu
 public | blogs                       | table             | ubuntu
 public | comment_counts_rank_by_blog | materialized view | ubuntu
 public | comments                    | table             | ubuntu
 public | members                     | table             | ubuntu
(5 rows)


SELECT *
FROM comment_counts_rank_by_blog;

EXPLAIN
SELECT *
FROM comment_counts_rank_by_blog;
```

``` sql
-- view
CREATE VIEW comment_counts_rank_by_blog_view AS
SELECT
    c.blog_code,
    COUNT(c) AS count,
    SUBSTRING(b.title FROM 1 FOR 30)  AS title,
    to_timestamp(b.timestamp::bigint / 1000000000) AS date,
    m.name AS name,
    b.link AS link
FROM comments AS c
INNER JOIN blogs AS b ON c.blog_code = b.code
INNER JOIN members AS m ON b.member_code = m.code
GROUP BY c.blog_code, b.timestamp, m.name, b.link, b.title
ORDER BY count DESC;


nogi-official=# \d
                           List of relations
 Schema |               Name               |       Type        | Owner  
--------+----------------------------------+-------------------+--------
 public | blog_images                      | table             | ubuntu
 public | blogs                            | table             | ubuntu
 public | comment_counts_rank_by_blog      | materialized view | ubuntu
 public | comment_counts_rank_by_blog_view | view              | ubuntu
 public | comments                         | table             | ubuntu
 public | members                          | table             | ubuntu
(6 rows)


EXPLAIN
SELECT *
FROM comment_counts_rank_by_blog_view;

-- view はクエリの保存のみで、取得時はクエリが実行されることが分かる。
nogi-official=# EXPLAIN
nogi-official-# SELECT *
nogi-official-# FROM comment_counts_rank_by_blog_view;
                                              QUERY PLAN                                              
------------------------------------------------------------------------------------------------------
 Subquery Scan on comment_counts_rank_by_blog_view  (cost=54.08..57.08 rows=240 width=728)
   ->  Sort  (cost=54.08..54.68 rows=240 width=1276)
         Sort Key: (count(c.*)) DESC
         ->  HashAggregate  (cost=38.59..44.59 rows=240 width=1276)
               Group Key: c.blog_code, b."timestamp", m.name, b.link, b.title
               ->  Hash Join  (cost=21.45..34.99 rows=240 width=1562)
                     Hash Cond: ((c.blog_code)::text = (b.code)::text)
                     ->  Seq Scan on comments c  (cost=0.00..12.40 rows=240 width=416)
                     ->  Hash  (cost=21.20..21.20 rows=20 width=1228)
                           ->  Hash Join  (cost=10.45..21.20 rows=20 width=1228)
                                 Hash Cond: ((m.code)::text = (b.member_code)::text)
                                 ->  Seq Scan on members m  (cost=0.00..10.40 rows=40 width=164)
                                 ->  Hash  (cost=10.20..10.20 rows=20 width=1228)
                                       ->  Seq Scan on blogs b  (cost=0.00..10.20 rows=20 width=1228)
(14 rows)
```

## how to refresh data

``` sql
REFRESH MATERIALIZED VIEW view_name;
```


``` sql
nogi-official=# SELECT *
FROM comment_counts_rank_by_blog;
 blog_code | count | title  |          date          |   name    |                link                 
-----------+-------+--------+------------------------+-----------+-------------------------------------
 1         |     3 | title  | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
 3         |     2 | title3 | 2023-05-14 11:50:26+00 | john doe2 | https://www.instagram.com/john.doe/
 2         |     1 | title2 | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
(3 rows)

nogi-official=# INSERT INTO comments(
nogi-official(#     code,
nogi-official(#     blog_code,
nogi-official(#     content,
nogi-official(#     name,
nogi-official(#     timestamp
nogi-official(# ) VALUES (
nogi-official(#     '7',
nogi-official(#     '1',
nogi-official(#     'content7',
nogi-official(#     'name5',
nogi-official(#     '1684065026000000000'
nogi-official(# );
INSERT 0 1

-- NOT updated
nogi-official=# SELECT *
FROM comment_counts_rank_by_blog;
 blog_code | count | title  |          date          |   name    |                link                 
-----------+-------+--------+------------------------+-----------+-------------------------------------
 1         |     3 | title  | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
 3         |     2 | title3 | 2023-05-14 11:50:26+00 | john doe2 | https://www.instagram.com/john.doe/
 2         |     1 | title2 | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
(3 rows)


nogi-official=# EXPLAIN REFRESH MATERIALIZED VIEW comment_counts_rank_by_blog;
                QUERY PLAN                 
-------------------------------------------
 Utility statements have no plan structure
(1 row)

-- REFRESH !!!!
nogi-official=# REFRESH MATERIALIZED VIEW comment_counts_rank_by_blog;
REFRESH MATERIALIZED VIEW

nogi-official=# SELECT *
FROM comment_counts_rank_by_blog;
 blog_code | count | title  |          date          |   name    |                link                 
-----------+-------+--------+------------------------+-----------+-------------------------------------
 1         |     4 | title  | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
 3         |     2 | title3 | 2023-05-14 11:50:26+00 | john doe2 | https://www.instagram.com/john.doe/
 2         |     1 | title2 | 2023-05-14 11:50:26+00 | john doe  | https://www.instagram.com/john.doe/
(3 rows)
```

