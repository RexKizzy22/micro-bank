setup:
	setenv LOCAL_DB_URL "$(LOCAL_DB_URL)"
	setenv DOCKER_DB_URL "$(DOCKER_DB_URL)" 
	setenv AWS_RDS_DB_URL "$(AWS_RDS_DB_URL)"
	setenv ECR_URL "$(ECR_URL)"
	setenv PG_CONTAINER_NAME "$(PG_CONTAINER_NAME)"
	setenv PG_IMAGE "$(PG_IMAGE)"
	setenv POSTGRES_PASSWORD "$(POSTGRES_PASSWORD)"
	setenv POSTGRES_USER "$(POSTGRES_USER)"
	setenv FILE_NAME "$(FILE_NAME)"

# start local server
server:
	# swag fmt
	go fmt ./...
	gow run main.go

# create a migration file
migration:
	migrate create -dir db/migration -ext sql -seq $(FILE_NAME)

# generate REST API docs from the available services' docstrings
swag:
	swag init

# generate a random string with 32 characters
randkey:
	openssl rand -hex 64 | head -c 32

# run postgres container using postgres official image
postgres:
	docker run --name=$(PG_CONTAINER_NAME) -p 5432:5432 \
		-e GIN_MODE=release -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_USER=$(POSTGRES_USER) -d postgres:15-alpine

# stop running container instance of the postgres image
stop-postgres:
	docker stop $(PG_CONTAINER_NAME)

# start stopped container instance of the postgres image
start-postgres:
	docker start $(PG_CONTAINER_NAME)

# remove stopped container instance of the postgres image
rm-postgres:
	docker rm $(PG_CONTAINER_NAME)

# Connect postgres container to microbank server using a common container network
postgres-dk:
	# docker run --name=postgres15 --network bank-network -p 5430:5432 \
		-e GIN_MODE=release -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_USER=$(POSTGRES_USER) -d postgres:15-alpine

# create a new database in postgres container
createdb:
	docker exec -it $(PG_CONTAINER_NAME) createdb --username=postgres --owner=postgres microbank

# delete the database in postgres container
dropdb:
	docker exec -it $(PG_CONTAINER_NAME) dropdb --username=postgres microbank

# run database queries in postgres container
querydb:
	docker exec -it $(PG_CONTAINER_NAME) psql -U postgres microbank

# migration with local database
migrateup:
	migrate -path db/migration -database $(LOCAL_DB_URL) -verbose up

migratedown:
	migrate -path db/migration -database $(LOCAL_DB_URL) -verbose down

migrateup1:
	migrate -path db/migration -database $(LOCAL_DB_URL) -verbose up 1

migratedown1:
	migrate -path db/migration -database $(LOCAL_DB_URL) -verbose down 1

# migration with docker compose
migrateup-dk:
	migrate -path db/migration -database $(DOCKER_DB_URL) -verbose up

migratedown-dk:
	migrate -path db/migration -database $(DOCKER_DB_URL) -verbose down

migrateup1-dk:
	migrate -path db/migration -database $(DOCKER_DB_URL) -verbose up 1

migratedown1-dk:
	migrate -path db/migration -database $(DOCKER_DB_URL) -verbose down 1

# migration with remote database on AWS RDS
migrateup-rmt:
	migrate -path db/migration -database $(AWS_RDS_DB_URL) -verbose up

migratedown-rmt:
	migrate -path db/migration -database $(AWS_RDS_DB_URL) -verbose down

migrateup1-rmt:
	migrate -path db/migration -database $(AWS_RDS_DB_URL) -verbose up 1

migratedown1-rmt:
	migrate -path db/migration -database $(AWS_RDS_DB_URL) -verbose down 1

# Generate models from sql queries and as well generate repositories to communicate with the database
sqlc:
	sqlc generate

# Run all test, output verbose logs and output coverage information
test:
	go test -v -cover -short ./...

# Generate database mock utilities for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Rexkizzy22/micro-bank/db/sqlc Store

# Generate database mock utilities for testing redis task
task-mock:
	mockgen -package mockwk -destination task/mock/distributor.go github.com/Rexkizzy22/micro-bank/task TaskDistributor

# Retrive authentication token from AWS ECR in order to gain access to remote container
awsecrlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin$(ECR_URL)

# Retrieve secrets from AWS Secret Manager and save them to app.env
awssecrets:
	aws secretsmanager get-secret-value --secret-id microbank -query SecretString \ 
		--output text | jq.'to_entries|map("\(.key)=\(.value)")|.[]' >> app.env

awscurrentuser:
	aws sts get-caller-identity

# Configure kubeconfig file to use AWS context
kubeconfig:
	aws eks update-kubeconfig --name microbank --region us-east-1

# Select a particular context for kubectl to use when there are multiple contexts registered
k8scontext:
# 	kubectl config use-context arn:aws:eks:us-west-2:335858042864:cluster/microbank 

# Compile the gRPC interface definition language and serve assets from a statik server
proto:
	rm -f pb/*.go
	rm -f gapi/docs/swagger/*.swagger.yaml
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=gapi/docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=microbank,output_format=yaml \
	proto/*.proto
	statik -src=./gapi/docs/swagger -dest=./gapi

# run gRPC calls using Evans client
evans:
	evans --host localhost --port 9090 -r repl

# generate database documentation using dbdiagram.io CLI
db_doc:
	dbdocs build db/microbank.dbml

# generate database schema from schema design tool by dbdiagram.io
db_schema:
	dbml2sql --postgres -o db/microbank.sql db/microbank.dbml

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-alpine

.PHONY: postgres createdb dropdb querydb migrateup migrateup1 migratedown migratedown1 migrateup-dock migrateup1-dock migratedown-dock migratedown1-dock migrateup-remote migrateup1-remote migratedown-remote migratedown1-remote sqlc test server swag proto evans db_doc db_schema awsecrlogin awssecrets awscurrentuser kubeconfig k8scontext
