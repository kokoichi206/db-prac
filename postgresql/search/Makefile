DB_URL=postgres://ubuntu:sakamichi@localhost:5435/psyllium?sslmode=disable

# FIXME: 各種ファイルに定義が書き込まれてるのをなんとかしたい。
SERVICE="postgres"
USER="ubuntu"
DB_NAME="nogi-official"

psql:	## docker compose で立てた DB 内に入る。
	docker compose exec $(SERVICE) psql -U $(USER) $(DB_NAME)

psql-test:	## docker compose で立てた DB 内に入る。
	docker compose exec $(SERVICE)test psql -U $(USER) $(DB_NAME)
