## View vs Materialized View

View は複雑なクエリを保存してるのみで、実行時はそのクエリが発行される。

Materialized View はクエリの結果が保存され、一種のキャッシュとして機能する。
データは静的なため手動または定期的に更新（REFRESH）が必要。

- [PostgreSQL View](https://www.postgresql.jp/document/9.1/html/sql-createview.html)
- [PostgreSQL M](https://www.postgresql.jp/docs/9.4/rules-materializedviews.html)
