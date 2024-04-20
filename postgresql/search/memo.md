## data 準備

``` sh
# table を絞る。
pg_dump --username=ubuntu -t members -t blogs nogi-official > backup_members_blogs
# 上のやつをリストア。
## -T: tty の割り当てを無効化。
dc exec -T postgres psql -U ubuntu -d nogi-official < backup_members_blogs

dc exec -T postgres psql -U ubuntu -d nogitest < backup_members_blogs

dc exec -T postgrestest psql -U ubuntu -d nogi-official < backup_members_blogs

```

## [document](https://www.postgresql.jp/docs/15/textsearch.html)

- intro
  - 辞書の活用
    - インデックスしたくないストップワードの定義
    - Ispell の活用
  - tsvector というデータ型？
  - @@ 演算子

``` sql
SELECT title
FROM pgweb
WHERE to_tsvector('english', body) @@ to_tsquery('english', 'friend');
```

``` sql
SELECT title FROM blogs WHERE to_tsvector('japanese', content) @@ to_tsquery('japanese', '本当に');

-- 検索自体はできてそう
SELECT link FROM blogs WHERE to_tsvector(content) @@ to_tsquery('本当に本当に');

SELECT link FROM blogs WHERE to_tsvector('english', content) @@ to_tsquery('english', '本当に本当に');

```

## pg_bigm 拡張

``` diff
+ shared_preload_libraries = 'pg_bigram'  # (change requires restart)
- shared_preload_libraries = ''  # (change requires restart)
```

``` sh
apt update
apt install -y postgresql-server-dev-15 make gcc wget libicu-dev
wget https://ja.osdn.net/dl/pgbigm/pg_bigm-1.2-20200228.tar.gz
tar zxf pg_bigm-1.2-20200228.tar.gz
cd pg_bigm-1.2-20200228 && make USE_PGXS=1 && make USE_PGXS=1 install
echo shared_preload_libraries='pg_bigm' >> /var/lib/postgresql/data/postgresql.conf


make psql
CREATE EXTENSION pg_bigm;
```


```sql
SELECT * FROM blogs WHERE content LIKE '%本当に本当に偉い%'```
```

## [pgroonga](https://pgroonga.github.io/ja/tutorial/)

- [PGroonga対pg_bigm](https://pgroonga.github.io/ja/reference/pgroonga-versus-pg-bigm.html)

[install](https://pgroonga.github.io/ja/install/ubuntu.html)

``` sh
# postgresql の DockerImage では Debian の bookworm
$ lsb_release -sc
No LSB modules are available.
bookworm
```

https://pgroonga.github.io/ja/install/debian.html

``` sh
apt update
apt install -y -V ca-certificates lsb-release wget
wget https://apache.jfrog.io/artifactory/arrow/$(lsb_release --id --short | tr 'A-Z' 'a-z')/apache-arrow-apt-source-latest-$(lsb_release --codename --short).deb
apt install -y -V ./apache-arrow-apt-source-latest-$(lsb_release --codename --short).deb
wget https://packages.groonga.org/debian/groonga-apt-source-latest-$(lsb_release --codename --short).deb
apt install -y -V ./groonga-apt-source-latest-$(lsb_release --codename --short).deb
apt update


echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release --codename --short)-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -


apt update
apt install -y -V postgresql-15-pgdg-pgroonga
```

チュートリアル

``` sql
CREATE EXTENSION IF NOT EXISTS pgroonga;

-- 検索用の Index をはる。
CREATE INDEX pgroonga_content_index ON blogs USING pgroonga (content);

SELECT title, link FROM blogs WHERE content &@ '本当に';

-- スコア用の Index をはる。
CREATE INDEX pgroonga_score_blogs_content_index
ON blogs
USING pgroonga (content);

SET enable_seqscan = off;
SELECT title, pgroonga_score(tableoid, ctid) AS score
  FROM blogs
 WHERE content &@ '本当に' OR content &@ '偉い';

SELECT m.name, b.link, pgroonga_score(b.tableoid, b.ctid) AS score
FROM blogs AS b
INNER JOIN members AS m ON b.member_code = m.code
WHERE content &@ '本当に' OR content &@ '偉い'
ORDER BY score DESC
;

SELECT m.name, b.link, pgroonga_score(blogs) AS score
FROM blogs AS b
INNER JOIN members AS m ON b.member_code = m.code
WHERE content &@ '本当に' OR content &@ '偉い'
ORDER BY score DESC
;
```

### スコアについて

[pgroonga_score 関数](https://pgroonga.github.io/ja/reference/functions/pgroonga-score.html)

> 現在のところ、スコアーの値は「何個キーワードが含まれていたか」（TF、Term Frequency）です。
> Groongaはどのようにスコアーを計算するかをカスタマイズすることができます。
> しかし、PGroongaはまだその機能をサポートしていません。

### ハイライト

``` sql
SELECT pgroonga_highlight_html(content,
    pgroonga_query_extract_keywords ('本当に本当に偉い')) AS highlighted_content,
    pgroonga_score(b) AS score
FROM blogs AS b
INNER JOIN members AS m ON b.member_code = m.code
WHERE
	CONTENT &@~ '本当に本当に偉い'
ORDER BY score DESC
;
```

[KWIC](https://ja.wikipedia.org/wiki/KWIC)

## Links

- [PostgreSQLで全文検索を実現するには](https://techblog.recochoku.jp/9920)
