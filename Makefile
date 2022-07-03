postgres:
	docker run --name=postgres14 -p 5430:5432 -e POSTGRES_PASSWORD=SimpleBank -e POSTGRES_USER=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres simple-bank

dropdb:
	docker exec -it postgres14 dropdb --username=postgres simple-bank

querydb:
	docker exec -it postgres14 psql -U postgres simple-bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5430/simple-bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5430/simple-bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5430/simple-bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5430/simple-bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Rexkizzy22/simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb querydb migrateup migrateup1 migratedown migratedown1 sqlc test server