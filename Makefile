server: swag
	swag fmt
	go fmt ./...
	go run main.go

swag: 
	swag init

postgres:
	docker run --name=postgres14 --network bank-network -p 5432:5432 -e GIN_MODE=release -e POSTGRES_PASSWORD=SimpleBank -e POSTGRES_USER=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres simple-bank

dropdb:
	docker exec -it postgres14 dropdb --username=postgres simple-bank

querydb:
	docker exec -it postgres14 psql -U postgres simple-bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose up
	migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose up 1
	migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose down
	migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose down 1
	migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Rexkizzy22/simple-bank/db/sqlc Store

awsecrlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin 335858042864.dkr.ecr.us-west-2.amazonaws.com

awssecrets:
	aws secretsmanager get-secret-value --secret-id simple_bank -query SecretString --output text | jq.'to_entries|map("\(.key)=\(.value)")|.[]' >> app.env

kubeconfig:
	aws eks update-kubeconfig --name simple-bank --region us-west-2

.PHONY: postgres createdb dropdb querydb migrateup migrateup1 migratedown migratedown1 sqlc test server swag awssecrets