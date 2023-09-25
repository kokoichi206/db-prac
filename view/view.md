``` sql
CREATE VIEW members_1st_generation
AS SELECT code, name, english_name, img
FROM members
WHERE cate = '1期生';
```

``` sh
\d members_1st_generation ;
                  View "public.members_1st_generation"
    Column    |          Type          | Collation | Nullable | Default 
--------------+------------------------+-----------+----------+---------
 code         | character varying(32)  |           |          | 
 name         | character varying(32)  |           |          | 
 english_name | character varying(32)  |           |          | 
 img          | character varying(255) |           |          | 
```

## insert

``` sql
INSERT INTO members_1st_generation
VALUES ('117117', 'John Doe', 'John Doe', 'pien.jpg');
```

なるほど、これは親のテーブルに挿入しようとしたことになるらしい！！

``` sh
nogi-official=# INSERT INTO members_1st_generation
nogi-official-# VALUES ('117117', 'John Doe', 'John Doe', 'pien.jpg');
ERROR:  null value in column "kana" of relation "members" violates not-null constraint
DETAIL:  Failing row contains (117117, John Doe, John Doe, null, null, pien.jpg, null, null, null, null, null, null, null, null, null).
```
