server: swag
	swag fmt
	go fmt ./...
	gow run main.go

swag: docs
	swag init

randkey:
	openssl rand -hex 64 | head -c 32

postgres:
	# Connect container to localhost
	docker run --name=postgres14 -p 5432:5432 -e GIN_MODE=release -e POSTGRES_PASSWORD=SimpleBank -e POSTGRES_USER=postgres -d postgres:14-alpine

	# Connect using a common container network
	# docker run --name=postgres14 --network bank-network -p 5432:5432 -e GIN_MODE=release -e POSTGRES_PASSWORD=SimpleBank -e POSTGRES_USER=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres simple-bank

dropdb:
	docker exec -it postgres14 dropdb --username=postgres simple-bank

querydb:
	docker exec -it postgres14 psql -U postgres simple-bank

migrateup:
	# Connect to the local database
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5432/simple-bank?sslmode=disable" -verbose up

	# Connect between the database container and the app container
	# migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose up

	# Connect to the remote database on AWS RDS
	# migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	# Connect to the local database
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5432/simple-bank?sslmode=disable" -verbose up 1

	# Connect between the database container and the app container
	# migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose up 1

	# Connect to the remote database on AWS RDS
	# migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown:
	# Connect to the local database
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@localhost:5432/simple-bank?sslmode=disable" -verbose down

	# Connect between the database container and the app container
	# migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose down

	# Connect to the remote database on AWS RDS
	# migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose down

migratedown1:
	# Connect to the local database
	migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose down 1

	# Connect between the database container and the app container
	# migrate -path db/migration -database "postgresql://postgres:SimpleBank@postgres14:5432/simple-bank?sslmode=disable" -verbose down 1

	# Connect to the remote database on AWS RDS
	# migrate -path db/migration -database "postgresql://postgres:Kizito22@simple-bank.cs5zwlono2zn.us-west-2.rds.amazonaws.com:5432/simple_bank" -verbose down 1

# Generate models from sql queries and as well generate repositories to communicate with the database
sqlc:
	sqlc generate

# Run all test and output coverage information
test:
	go test -v -cover ./...

# Generate database mock utilities for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Rexkizzy22/simple-bank/db/sqlc Store

# Retrive authentication token from AWS ECR in order to gain access to remote container
awsecrlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin 335858042864.dkr.ecr.us-west-2.amazonaws.com

# Retrieve secrets from AWS Secret Manager and save them to app.env
awssecrets:
	aws secretsmanager get-secret-value --secret-id simple_bank -query SecretString --output text | jq.'to_entries|map("\(.key)=\(.value)")|.[]' >> app.env

# Configure kubeconfig file to use AWS context
kubeconfig:
	aws eks update-kubeconfig --name simple-bank --region us-west-2

# Select a particular context for kubectl to use when there are multiple contexts registered
k8scontext:
	kubectl config use-context arn:aws:eks:us-west-2:335858042864:cluster/simple-bank

.PHONY: postgres createdb dropdb querydb migrateup migrateup1 migratedown migratedown1 sqlc test server swag awssecrets