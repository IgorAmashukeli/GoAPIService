.PHONY: test app

test:
	docker compose up --build --abort-on-container-exit test

app:
	docker compose up --build app

reset-db:
	docker compose exec -T db psql -U user -d mydb -c "DROP TABLE IF EXISTS answers CASCADE; DROP TABLE IF EXISTS questions CASCADE; DROP TABLE IF EXISTS goose_db_version;"