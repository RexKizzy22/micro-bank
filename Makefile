setup:
	setenv LOCAL_DB_URL "$(LOCAL_DB_URL)"
	setenv DOCKER_DB_URL "$(DOCKER_DB_URL)" 
	setenv AWS_RDS_DB_URL "$(AWS_RDS_DB_URL)"
	setenv ECR_URL "$(ECR_URL)"

# start local server
server: swag
	swag fmt
	go fmt ./...
	gow run main.go

# create a migration file
migration:
	migrate create -dir db/migration -ext sql -seq <filename>

# generate REST API docs from the available services' docstrings
swag:
	swag init

randkey:
	openssl rand -hex 64 | head -c 32

# run postgres container using postgres official image
postgres:
	docker run --name=postgres14 -p 5430:5432 \
		-e GIN_MODE=release -e POSTGRES_PASSWORD=MicroBank -e POSTGRES_USER=postgres -d postgres:14-alpine

# stop running container instance of the postgres image
stop-postgres:
	docker stop postgres14

# remove stopped container instance of the postgres image
rm-postgres:
	docker rm postgres14

# Connect postgres container to MicroBank server using a common container network
postgres-dock:
	# docker run --name=postgres14 --network bank-network -p 5432:5432 \
		-e GIN_MODE=release -e POSTGRES_PASSWORD=MicroBank -e POSTGRES_USER=postgres -d postgres:14-alpine

# create a new database in postgres container
createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres microbank

# delete the database in postgres container
dropdb:
	docker exec -it postgres14 dropdb --username=postgres microbank

# run database queries in postgres container
querydb:
	docker exec -it postgres14 psql -U postgres microbank

# Connect to the local database
migrateup:
	migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DOCKER_DB_URL)" -verbose down 1

# Connect between the database container and the app container
migrateup-dock:
	migrate -path db/migration -database "$(DOCKER_DB_URL)" -verbose up

migrateup1-dock:
	migrate -path db/migration -database "$(DOCKER_DB_URL)" -verbose up 1

migratedown-dock:
	migrate -path db/migration -database "$(DOCKER_DB_URL)" -verbose down

migratedown1-dock:
	migrate -path db/migration -database "$(DOCKER_DB_URL)" -verbose down 1

# Connect to the remote database on AWS RDS
migrateup-remote:
	migrate -path db/migration -database "$(AWS_RDS_DB_URL)" -verbose up

migrateup1-remote:
	migrate -path db/migration -database "$(AWS_RDS_DB_URL)" -verbose up 1

migratedown-remote:
	migrate -path db/migration -database "$(AWS_RDS_DB_URL)" -verbose down

migratedown1-remote:
	migrate -path db/migration -database "$(AWS_RDS_DB_URL)" -verbose down 1

# Generate models from sql queries and as well generate repositories to communicate with the database
sqlc:
	sqlc generate

# Run all test, output verbose logs and output coverage information
test:
	go test -v -cover ./...

# Generate database mock utilities for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Rexkizzy22/micro-bank/db/sqlc Store

# Retrive authentication token from AWS ECR in order to gain access to remote container
awsecrlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin "$(ECR_URL)"

# Retrieve secrets from AWS Secret Manager and save them to app.env
awssecrets:
	aws secretsmanager get-secret-value --secret-id microbank -query SecretString \ 
		--output text | jq.'to_entries|map("\(.key)=\(.value)")|.[]' >> app.env

awscurrentuser:
	aws sts get-caller-identity

awscurrentuser:
	aws sts get-caller-identity

# Configure kubeconfig file to use AWS context
kubeconfig:
	aws eks update-kubeconfig --name microbank --region us-east-1

# Select a particular context for kubectl to use when there are multiple contexts registered
k8scontext:
# 	kubectl config use-context arn:aws:eks:us-west-2:335858042864:cluster/microbank 

# Compile the GRPC interface definition language and serve assets from a statik server
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=gapi/swagger --openapiv2_opt=allow_merge=true,merge_file_name=microbank,output_format=yaml \
	proto/*.proto
	statik -src=./gapi/swagger -dest=./gapi

# run GRPC calls using Evans client
evans:
	evans --host localhost --port 9090 -r repl

# generate database documentation using dbdiagram.io CLI
db_doc:
	dbdocs build dbdoc/microbank.dbml

# generate database schema from schema design tool by dbdiagram.io
db_schema:
	dbml2sql --postgres -o dbdoc/microbank.sql dbdoc/microbank.dbml

.PHONY: postgres createdb dropdb querydb migrateup migrateup1 migratedown migratedown1 migrateup-dock migrateup1-dock migratedown-dock migratedown1-dock migrateup-remote migrateup1-remote migratedown-remote migratedown1-remote sqlc test server swag proto evans db_doc db_schema awsecrlogin awssecrets awscurrentuser kubeconfig k8scontext