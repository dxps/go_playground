# .PHONY: db-postgres db-migrate
.PHONY: db-migrate

# Replaced by run_db.sh script
# db-postgres:
# 	docker run --name go-reddit-db --rm -ti -p 54326:5432 -e POSTGRES_PASSWORD=secret postgres

db-migrate:
	go_migrate -source file://db/migrations \
	           -database postgres://go-reddit:secret@localhost:54326/go-reddit?sslmode=disable \
			   up

