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
        - 総当たり方式
        - 扱うテーブル数が閾値（通常12）を超えた場合、遺伝的問い合わせ最適化
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


## part 2

### sec 5
理想的なテーブル設計は RDBMS に依存しないことだが、現実には各 RDBMS に特化した知識が必要となることもある。

#### データ型
- 文字型
- 数値データ型
- 日付/時刻データ型
- バイナリ列データ型
    - 任意の値のバイトを格納できる

#### 文字型
- text
    - 基本的には text 型を用いることが[推奨されている](https://www.postgresql.jp/document/9.5/html/datatype-character.html)
    - varchar は格納時のサイズチェックが行われるぶん、わずかに text より遅い

#### 日付/時刻データ型
時刻型にタイムゾーンを含めるかどうかは、システムで扱う情報が時間帯をまたがる必要があるかどうかによる。１つの時間帯に時刻情報が収まるようなケースでは（日本国内の自国の身を意識すればよいケース）では、タイムゾーンなしの型で問題ないと考えられる！

#### バイナリ列データ型
基本的に `bytea` 型を使用する。格納するデータが非常に大きい場合等は、ラージオブジェクトの使用も検討する。bytea 型では最大 1GB で、ラージオブジェクトでは最大 4TB。

性能の観点では、格納するデータ量が大きい場合、内部でのデータコピー量が多くなることから bytea 型へのアクセス性能が悪くなる傾向がある。基準として、**格納するデータが 100KB までは bytea 型を、それ以上の場合にはラージオブジェクトを使用する！**

[ラージオブジェクト](https://www.postgresql.jp/document/9.5/html/largeobjects.html)は、管理方法が全く異なる。データ型があるわけではなく、ラージオブジェクト識別子の列ができる形。


#### 制約
PosgtreSQL では主キーに対して、暗黙のうちに btree インデックスが設定される。

一意性制約に対しても、暗黙の B-tree インデックスがはられる。

外部キーでの注意点

- 外部キーには暗黙的なインデックスが設定されない
- テーブル間で外部キーの型を一致させる！

#### TOAST
過大属性格納技法（The Oversized-Attribute Storage Technique）は、非常に大きな列の値を格納する実装技法。

TOAST 化はテーブルに格納しようとする行のサイズが 2KB を超えるときに実行される。通常で十分に効率化されているため、TOAST を発生させないように意識する必要は特にない。

#### 結合を意識したテーブル設計
テーブルの結合処理、JOIN、は性能に大きな影響を与える！

結合数が非常に多くなると、通常の問い合わせ最適化より若干精度の劣る問い合わせ最適化が行われることがあり、必要に応じて**正規化を崩して結合数を減らすことも検討**する。

### sec 6
各ファイルをどこに、どのように格納するかといった物理的な面を考慮した設計。

PostgreSQL では、テーブル用のファイルもインデックス用のファイルもファイル形式は統一されている。
基本的には 8192 バイトのページと呼ばれる固定超領域が連続して配置されたもの。

#### テーブルファイルに対するアクセス
大別すると「シーケンシャルアクセス」と「インデックスアクセス」に分類される

#### WAL ファイルとアーカイブファイル
WAL ファイルは、先行書き込みログ（Write Ahead Logging）が格納される非常に重要なファイル。

#### HOT と FILLFACTOR
HOT（Heap Only Tuple）により、PostgreSQL の更新性能が大幅に向上した。（[仕組みについての記事](https://lets.postgresql.jp/documents/tutorial/hot_2)）

- UPDATE 時のインデックスエントリの追加処理をスキップする
- VACUUM 処理を待つことなく不要領域を再利用可能にする

HOT が有効になるのは、インデックスを持たない列への更新で、さらに更新対象の行と同じページ内に空きがあり、新しい行を挿入可能な場合。なお、次のような更新処理の場合、HOT は働かない。

HOT 機能自体は、ユーザが意図的に制御できないが、効果的に活用するために FILLFACTOR による物理設計を考慮する必要がある。[更新の多い処理は FILLFACTOR で空き領域を作成しておく](https://lets.postgresql.jp/documents/tutorial/hot3)。

#### Index
値の分布に偏りがある場合には、部分インデックスをはることも検討する（インデックスサイズの減少に効果的）

``` sql
-- 99% が 1,2 で残り 1% が 100 以上の場合、以下のようなインデックスが有効
CREATE INDEX test_idx ON test (value) WHERE value >= 100;
```


### sec 7
バックアップ計画を立てる際、リカバリ要件を明確にする！

- どの時点までのデータが必要か？
- どのくらいの時間で復旧させるか？


### sec 8
監視計画

#### 監視とは
「健全に動作している」ことを確認するための「監視項目の選定」「確認する間隔」「何をもって異常と判断するのかの閾値」を決定する。

#### 監視項目の選定
データベースが健全に動作しているか

- サーバーに問題がないか？
- PostgreSQL に問題がないか？

デフォルトで PostgreSQL ログは標準出力に出力され、後から確認ができない。このため、きちんとログをファイルに書き出す設定を行う必要あり。


``` conf
# postgresql.conf
log_destination = 'csvlog'
logging_collector = on
log_directory = '/var/lib/postgresql/data/'
log_filename = 'test.log'
log_file_mode = 0777
log_min_messages = debug5
log_line_prefix = '(user:%u access to database:%d at [%m], pid: %p)'
```

ログが極端に大きくなりすぎると、保守に手間がかかったり、いざ参照したい時には時間がかかったりとよいことがない。最悪のケースではログが肥大化して DB が停止、、、ということもある。そのため、適度なタイミングでログのローテーションを行う設定が必要！（ログローテは DB に限らない）


#### サーバ設定
PostgreSQL はクライアントからの要求を１つのプロセスが処理するプロセスモデルのアーキテクチャを採用している。

#### メモリの設定
HDD などからデータを取り出す時間と、メモリからデータを取り出す時間は、数100~数10万倍の性能差があると言われている！

メモリを生かすため、データアクセス時にデータベースファイルをページ単位でメモリ上に展開し、繰り返しアクセスする場合の処理性能を高めている。

PostgreSQL のメモリ設定の中でも特に重要なパラメータは「shared_buffers」であり、PostgreSQL が共有バッファのために確保する共有メモリのサイズを設定する。

#### ディスクの設定
一般的にデータベースは、ディスク性能がシステムのボトルネックになりやすい傾向がある。

#### OSのディスク設定
OS 設定のディスクに関連する設定では、I/O スケジューラの設定が有効。I/O スケジューラは、OS 上で動作しているさまざまなプロセスからの I/O 要求をどのように処理するかを定めているパラメータ。

``` sh
cat /sys/block/sda/queue/scheduler
cat /sys/block/vda/queue/scheduler
```

- noop
    - OS はスケジューラに関与しない
- anticipatory
    - I/O 要求に対して HDD ドライブの中の物理的な配置が近いデータを優先して処理する
- deadline
    - I/O 要求の待ち時間に限界値を設け、限界値に近い I/O を優先して処理する
- cfq
    - I/O 要求すべてを均等に処理する
    - CentOS のデフォルト
    - プロセスから小さい I/O 要求が発生する場合に適した設定

``` sh
# ubuntu22
$ cat /sys/block/sda/queue/scheduler
[mq-deadline] none

# Alpine Linux
$ cat /sys/block/vda/queue/scheduler 
[none] mq-deadline kyber bfq
```

PostgreSQL では deadline を推奨


