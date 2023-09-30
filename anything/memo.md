## TILE

``` sql
SELECT m.name, COUNT(*) , NTILE(5) OVER (ORDER BY COUNT(*) DESC) FROM blogs INNER JOIN members AS m ON m.code = blogs.member_code GROUP BY m.name;
```

ちょうどで割れなかった場合は、上から1つずつ増えていくらしい。

``` sql
SELECT rank, COUNT(*) FROM (
    SELECT m.name, NTILE(5) OVER (ORDER BY COUNT(*) DESC) AS rank
        FROM blogs INNER JOIN members AS m ON m.code = blogs.member_code
        GROUP BY m.name
) AS ranks GROUP BY rank;

nogi-official(# ) AS ranks GROUP BY rank;
 rank | count 
------+-------
    1 |    10
    2 |    10
    3 |     9
    4 |     9
    5 |     9
(5 rows)
```

## Domain

- [postgresql create domain](https://www.postgresql.jp/docs/9.2/sql-createdomain.html)
- [postgresql データ型](https://www.postgresql.jp/docs/9.2/datatype.html)

``` sql
CREATE DOMAIN domain_name AS data_type
[DEFAULT default_value]
[CONSTRAINT constraint_name CHECK (condition)];
```

- データの整合性を保つための制約を一元管理できる
- 同じ制約を持つカラムを複数のテーブルで使用する際に便利
- 複雑な制約を多用すると、データの挿入や更新の際にパフォーマンスの問題が生じる可能性あり

例

``` sql
-- 正の整数のみを許可するドメインを作成
CREATE DOMAIN positive_integer AS INTEGER
CHECK (VALUE > 0);

-- ドメインを使用してテーブルを作成
CREATE TABLE products (
    product_id positive_integer PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL
);

DROP DOMAIN IF EXISTS domain_name;
```

``` sql
-- 電話番号ドメインの作成
CREATE DOMAIN phone_number AS VARCHAR(15)
CHECK (
    VALUE ~ '^\+?[0-9\-]+$' AND
    LENGTH(VALUE) >= 10 AND LENGTH(VALUE) <= 15
);

CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone phone_number
);
CREATE TABLE employees (
    employee_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone phone_number
);
```

``` sql
CREATE DOMAIN url AS VARCHAR(255)
CHECK (
    VALUE ~ '^https://$' AND
    LENGTH(VALUE) >= 10 AND LENGTH(VALUE) <= 15
);

CREATE DOMAIN graduation AS VARCHAR(4)
CHECK (VALUE IN ('YES', 'NO'));

-- ALTER TABLE テーブル名 ALTER COLUMN カラム名 TYPE データ型
ALTER TABLE members ALTER COLUMN graduation TYPE graduation;
```

定義されたドメインの確認

``` sql
nogi-official=# \dD
                                                                                   List of domains
 Schema |    Name    |          Type          | Collation | Nullable | Default |                                                Check                                                 
--------+------------+------------------------+-----------+----------+---------+------------------------------------------------------------------------------------------------------
 public | graduation | character varying(4)   |           |          |         | CHECK (VALUE::text = ANY (ARRAY['YES'::character varying, 'NO'::character varying]::text[]))
 public | url        | character varying(255) |           |          |         | CHECK (VALUE::text ~ '^https://$'::text AND length(VALUE::text) >= 10 AND length(VALUE::text) <= 15)
(2 rows)
```


``` sql
-- スキーマ一覧？
SELECT *
FROM information_schema.tables;

-- スキーマ内のテーブル一覧確認
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'pg_catalog';

           table_name            
---------------------------------
 pg_statistic
 pg_type
 pg_foreign_table
 pg_authid
 pg_shadow
 pg_roles
 pg_statistic_ext_data
 pg_settings
 pg_file_settings
 pg_hba_file_rules
 pg_ident_file_mappings
 pg_config
 pg_shmem_allocations
 pg_backend_memory_contexts
 pg_user_mapping
 pg_stat_activity
 pg_replication_origin_status
 pg_subscription
 pg_attribute
 pg_proc
 pg_class
 pg_attrdef
 pg_constraint
 pg_inherits
 pg_index
 pg_stat_replication
 pg_stat_slru
 pg_stat_wal_receiver
 pg_stat_recovery_prefetch
 pg_operator
 pg_opfamily
 pg_opclass
 pg_am
 pg_amop
 pg_amproc
 pg_language
 pg_largeobject_metadata
 pg_aggregate
 pg_statistic_ext
 pg_rewrite
 pg_trigger
 pg_event_trigger
 pg_description
 pg_cast
 pg_enum
 pg_namespace
 pg_conversion
 pg_depend
 pg_database
 pg_db_role_setting
 pg_tablespace
 pg_auth_members
 pg_shdepend
 pg_shdescription
 pg_ts_config
 pg_ts_config_map
 pg_ts_dict
 pg_ts_parser
 pg_ts_template
 pg_extension
 pg_foreign_data_wrapper
 pg_foreign_server
 pg_policy
 pg_replication_origin
 pg_default_acl
 pg_init_privs
 pg_seclabel
 pg_shseclabel
 pg_collation
 pg_parameter_acl
 pg_partitioned_table
 pg_range
 pg_transform
 pg_sequence
 pg_publication
 pg_publication_namespace
 pg_publication_rel
 pg_subscription_rel
 pg_group
 pg_user
 pg_policies
 pg_rules
 pg_views
 pg_tables
 pg_matviews
 pg_indexes
 pg_sequences
 pg_stats
 pg_stats_ext
 pg_stats_ext_exprs
 pg_publication_tables
 pg_locks
 pg_cursors
 pg_available_extensions
 pg_available_extension_versions
 pg_prepared_xacts
 pg_prepared_statements
 pg_seclabels
 pg_timezone_abbrevs
 pg_timezone_names
 pg_stat_sys_tables
 pg_stat_xact_sys_tables
 pg_stat_user_tables
 pg_stat_all_tables
 pg_stat_xact_all_tables
 pg_stat_xact_user_tables
 pg_statio_all_tables
 pg_statio_sys_tables
 pg_statio_user_tables
 pg_stat_all_indexes
 pg_stat_sys_indexes
 pg_stat_user_indexes
 pg_statio_all_indexes
 pg_statio_sys_indexes
 pg_statio_user_indexes
 pg_statio_all_sequences
 pg_statio_sys_sequences
 pg_statio_user_sequences
 pg_stat_subscription
 pg_stat_ssl
 pg_stat_gssapi
 pg_replication_slots
 pg_stat_replication_slots
 pg_stat_database
 pg_stat_database_conflicts
 pg_stat_user_functions
 pg_stat_xact_user_functions
 pg_stat_archiver
 pg_stat_bgwriter
 pg_stat_wal
 pg_stat_progress_analyze
 pg_stat_progress_vacuum
 pg_stat_progress_cluster
 pg_stat_progress_create_index
 pg_stat_progress_basebackup
 pg_stat_progress_copy
 pg_user_mappings
 pg_stat_subscription_stats
 pg_largeobject
(139 rows)



SELECT * FROM pg_catalog.pg_type;
\d pg_catalog.pg_type
```


## TRIGGER

``` sql
CREATE TABLE IF NOT EXISTS users (
    name varchar(32),
    bio text,
    CONSTRAINT users_pk PRIMARY KEY (name)
);

CREATE MATERIALIZED VIEW IF NOT EXISTS users_stat
AS
SELECT name, LENGTH(bio) FROM users;

-- うーん
CREATE OR REPLACE FUNCTION update_users_stat()
RETURNS TRIGGER AS $$
BEGIN
    -- users_statを更新するロジック
    REFRESH MATERIALIZED VIEW users_stat;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_users_stat
AFTER INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_stat();

INSERT INTO users (name, bio) VALUES (
    'bar',
    'from golang'
);
```

``` sql
nogi-official=# SELECT 1, COALESCE(null, 1, null, '3');
 ?column? | coalesce 
----------+----------
        1 |        1
(1 row)

SELECT name, CASE WHEN length > 12 THEN 'long' ELSE 'short' END AS result, length AS bio_length FROM users_stat;
```


## buffer cache

[サーバーのメモリの 1/4-1/2 にするといいらしい](https://thinkit.co.jp/cert/tech/10/1/3.htm)

``` sql
nogi-official=# CREATE EXTENSION pg_buffercache;
CREATE EXTENSION
nogi-official=# 
nogi-official=# SELECT count(*) FROM pg_buffercache;
 count 
-------
 16384
(1 row)
```
