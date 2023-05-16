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
❯ dc exec postgres psql -U ubuntu nogi-official
create database "nogi-official-member-only-test"

# このコマンドでは db があることを前提としてるみたい？
dc exec postgres psql -username=ubuntu nogi-official-member-only-test < backup/backup_members

dc exec postgres -U ubuntu 

```


