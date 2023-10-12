.DEFAULT_GOAL := serve

serve:
	go run main.go

# FIXME: 各種ファイルに定義が書き込まれてるのをなんとかしたい。
SERVICE="postgres"
USER="ubuntu"
DB_NAME="posgre"

psql:	## docker compose で立てた DB 内に入る。
	docker compose exec $(SERVICE) psql -U $(USER) $(DB_NAME)

db-up:
	docker compose up

db-down:
	docker compose down

help:	## https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
