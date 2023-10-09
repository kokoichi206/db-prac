.PHONY: run
run:
	go run main.go

SERVICE="postgres"
USER="ubuntu"
DB_NAME="postgres"

.PHONY: db-up
db-up:
	docker compose up

.PHONY: psql
psql:	## docker compose で立てた DB 内に入る。
	docker compose exec $(SERVICE) psql -U $(USER) $(DB_NAME)

.PHONY: logs
logs:	## docker compose で立てた DB のログを表示する。
	docker compose logs -f $(SERVICE)
