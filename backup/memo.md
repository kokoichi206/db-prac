## テキスト形式

``` sh
pg_dump -Fc nogi-official > backup_0516 

# 注: ここだけ psql のコマンド
psql nogi-official < backup_file名

psql nogi-official-member-only-test < backup_members
```


## アーカイブ形式

``` sh
dc exec postgres pg_dump -U ubuntu -d nogi-official -F c > backup_0516


pg_restore -C -d <db-name> backup_0516


pg_dump --username=ubuntu -t members -t blogs nogi-official > backup_members_blogs

dc exec postgres pg_dump --username=ubuntu -t members -t blogs nogi-official > backup_members_blogs
```


## restore

``` sh
dc exec postgres psql -U ubuntu nogi-official
create database "nogi-official-member-only-test";

# このコマンドでは db があることを前提としてるみたい？
dc exec postgres psql -U ubuntu -d nogi-official-member-only-test < backup/backup_members

dc exec postgres -U ubuntu 

```


## dc などを使う時の注意

Docker Compose（dc）のコンテナ内でコマンドを実行する場合、< など正しく認識されないことがある。
これはシェルがリダイレクトを処理するため、Docker Compose がリダイレクトを見る前にシェルが先に処理してしまうから。

``` sh
dc exec postgres bash -c "psql -U ubuntu -d nogi-official-member-only-test < /backup/data/backup_members"
```

ただし、この方法を使用する場合、SQLファイルがコンテナ内に存在している必要がある。
SQLファイルがホストマシンに存在する場合は、適切なボリュームマウントを使用してコンテナ内にファイルを利用できるようにする。

