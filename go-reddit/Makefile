.PHONY: postgres db-migrate

postgres:
	docker run --name go reddit-db --rm -ti -p 54326:5432 -e POSTGRES_PASSWORD=secret postgres

db-migrate:
	go_migrate -source file://db/migrations \
	           -database postgres://postgres:secret@localhost:54326/postgres?sslmode=disable \
			   up

