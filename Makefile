postgres:
	docker run --name postgres-latest -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d -p 5432:5432 postgres

createdb:
	docker exec -it postgres-latest createdb --username=root --owner=root $(POSTGRES_DATABASE)
	docker exec -it postgres-latest createdb --username=root --owner=root $(POSTGRES_DATABASE_DEV)

dropdb:
	docker exec -it postgres-latest dropdb $(POSTGRES_DATABASE)
	docker exec -it postgres-latest dropdb $(POSTGRES_DATABASE_DEV)

migrateup:
	migrate -path migration -database $(POSTGRES_URI) --verbose up
	migrate -path migration -database $(POSTGRES_URI_DEV) --verbose up

migratedown:
	migrate -path migration -database $(POSTGRES_URI) --verbose down
	migrate -path migration -database $(POSTGRES_URI_DEV) --verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...