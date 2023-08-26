run-postgres:
	docker run --name bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=test_local -d postgres:13-alpine

exec-postgres:
	docker exec -it bank psql -U root

create-db:
	docker exec -it bank createdb --username=root --owner=root bank

drop-db:
	docker exec -it bank dropdb bank

migrate-up:
	migrate -path db/migration -database "postgresql://root:test_local@localhost:5432/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:test_local@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

format:
	find $(PWD) -name "*.go" -exec gofmt -w {} \;

start:
	go run main.go

.PHONY: test
