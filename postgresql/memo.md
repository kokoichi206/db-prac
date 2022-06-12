## part 1
基本とは"簡単なこと"ではない

### sec 1
リレーショナルデータベースとオブジェクトデータベースの双方の能力を兼ね備えたオブジェクトリレーショナルデータベースに分類される。

オブジェクトデータベース機能として、ユーザ定義によりさまざまな機能拡張が可能！

- BSD ライセンスに類似した PostgreSQL ライセンス

### sec 2
RDBMS は、クエリの処理だけでなく、バッファの管理、ストレージへの書き込み制御、統計情報の収集など様々な制御を行なっている。PostgreSQL では複数のプロセスを動作させることで、複雑な制御を可能としている。

#### プロセス
``` sh
$ ps | grep postgres
    1 postgres  0:28 postgres
   51 postgres  0:01 postgres: checkpointer   
   52 postgres  0:10 postgres: background writer   
   53 postgres  0:10 postgres: walwriter   
   54 postgres  0:16 postgres: autovacuum launcher   
   55 postgres  0:32 postgres: stats collector   
   56 postgres  0:02 postgres: logical replication launcher   
   64 root      0:00 psql -h localhost -p 5432 -U postgres
 1519 postgres  0:00 postgres: postgres test_db 127.0.0.1(56634) idle
 5794 root      0:00 grep postgres
```

- マスタサーバプロセス
- ライタプロセス
    - 共有バッファ内の更新されたページを、対応するデータファイルに書き出す
- WAL ライタプロセス
    - Write Ahead Logging
    - WAL をディスクに書き出すプロセス
- チェックポインタプロセス
    - チェックポイント（全てのダーティページをデータファイルに反映し、特殊なチェックポイントレコードがログファイルに書き込まれた状態）を設定に従い、自動的に実行するプロセス
- 自動バキュームランチャと自動バキュームわーかプロセス
    - ランチャは設定に従ってわーかを起動
    - ワーカはテーブルに対して自動的にバキュームとアナライズを実行
- 統計情報コレクタプロセス
    - データベースの活動状況に関する稼働統計情報を一定間隔で収集するプロセス
- バックエンドプロセス  
    - クライアントから接続要求を受けたときに生成されるプロセス
    - クエリの実行

#### メモリ
メモリは、PostgreSQL サーバプロセス全体で共有される共有メモリ域と、バックエンドプロセスで確保されるプロセスメモリ域の２つに区別される！

#### ファイル
PosgtreSQL で使われるファイルの多くは、データベースクラスタと呼ばれるディレクトリは以下に生成される。initdb コマンドによって、OS ユーザにのみアクセス可能な権限（700）で作成される。

主なファイル

- PG_VERSION
- テーブルファイル
    - 複数の 8192 バイトの「ページ」によって構成される
- インデックスファイル
- TOAST ファイル
    - テーブル内に巨大な行を格納する場合に生成される特殊ファイル
- Free Space Map ファイル
- Visibility Map ファイル
- WAL ファイル
    - 更新操作を記録するファイル

### sec 3
設定ファイルの種類

- postgresql.conf
- pg_hba.conf
- pg_ident.conf
- recovery.conf
- pg_service.conf

``` sh
$ psql postgres -c "SHOW shared_buffers"
$ psql postgres -c "SHOW ALL"
```

``` postgresql
postgresql=# SET enable_seqscan = off;
# SET 文で変更可能な設定項目
postgresql=# SELECT name, context FROM pg_settings WHERE context IN ('user','superuser');
```

IPv4 での「0.0.0.0/0」、ipv6 での「::/0」は、それぞれ全ての IP アドレスを示す特殊な記法。CIDR（Classless Inter-Domain Routing）マスク。


### sec 4
処理/制御の基本

DBMS として必要な動作を、複数のサーバープロセスとプロセス間で共有するリソースによって制御している。

- マスタサーバプロセス
    - 他の全てのプロセスの親プロセス
    - バックエンドプロセスもここから fork() で起動される
- チェックポインタ
    - チェックポイントは、PostgreSQL がクラッシュした際にどの場所からリカバリ処理を行うかを示す
- バックエンドプロセス
    - 実際に SQL が実行されるプロセス

#### 問い合わせの実行
1. パーサ
    - 字句解析(flex)と構文解析(bison)
    - 存在するかの確認
    - SQL -> 問い合わせツリー
2. リライタ
    - ビューなどは、ルールを使って定義されている
    - 問い合わせツリー + ルール定義 -> 問い合わせツリー
3. プランナ/オプティマイザ
    - 作成後の実行計画は EXPLAIN 文で確認できる
    - 問い合わせツリー + 統計情報 -> 実行計画
4. エグゼキュータ
    - DML(Data Manipulation Language)のみを対象に処理を行う

- 実行計画の作成
    - 個々のテーブルに対するアクセス方法の選択
    - 結合方法の選択
- SQL の種別
    - DDL(Data Definition Language)
        - CREATE
        - DROP
        - ALTER ...
    - DML(Data Manipulation Language)
        - SELECT
        - INSERT ...
    - DCL(Data Control Language)
        - BEGIN
        - COMMIT
        - ROLLBACK

#### トランザクション
互いに関連する複数の処理を、トランザクションと呼ばれる不可分な処理単位として扱う。RDBMS の根幹となる機構！

トランザクションが満たすべき４つの要件

- 原子性（atomicity）
    - 複数の処理をまとめて、それらの全てが実行されたか、または全く実行されないかのどちらかの結果となること
- 一貫性（consistency）
    - 開始及び終了時点で、業務として規定された整合性を満たすこと
- 独立性（isolation）
    - 作業中のトランザクションによる更新は、確定するまで他のトランザクションから不可視になること
- 永続性（durability）
    - 確定したトランザクションの結果は永続的に保存されること

#### ロック
テーブル/行に対して明示的なロックを獲得できる。
テーブル単位のロックは LOCK 文を使用する。行単位のロックは SELECT 文のオプション指定として SELECT FOR UPDATE または SELECT FOR SHARE を使用する。

PostgreSQL にはデッドロックを検知する機構があり、検知した場合には起因となったトランザクションの片方を中断する！

#### 同時実行制御
PostgreSQL は追記型アーキテクチャを採用することで、MVCC（Multi Version Concurrency Control）と呼ばれる同時実行制御方式を実現している。
追記型アーキテクチャは、データの更新時に元々あったデータを直接更新するのではなく、更新前のデータはそのままに更新後のデータを追記する、という仕組み。




