# Microbank

**A minimalist implementation of a banking server, exposing endpoints for the following purposes:**

- creating a banking customer
- creating an account for a particular customer using one of the accepted currencies
- transfering money between customer accounts
- retrieving all accounts created by the customer in the current session
- retrieving one account created by the customer in the current session

**Read [NOTE.md](NOTE.md) for detailed information about this project**

## Run Microbank

You can run Microbank in several different ways:

### 1. Docker Compose - (The easiest)
- Download [Docker](https://www.docker.com/products/docker-desktop/)
- Clone this repository
- Comment out the environment variables for localhost and uncomment the ones for docker compose environment in **app.env**
- Start services
```bash
docker-compose up --build
```

### 2. Localhost (Macbook) - (Requires more installations on your machine)

**Microbank implements HTTP, gRPC Gateway and gRPC servers**

**gRPC Gateway server serves both HTTP and gRPC requests**

#### Pre-requisites
1. [Homebrew](https://brew.sh/)
2. PostgreSQL. Run `brew install postgresql`
3. [Go](https://go.dev/dl/)
4. golang-migrate. Run `brew install golang-migrate`
5. Install [redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/)

#### Run HTTP server
- Ensure these lines of code are commented out in the **main** function
```go
go runGatewayServer(config, store)
runGrpcServer(config, store)
```
- Comment out the environment variables for docker environment and uncomment the ones for localhost in **app.env** file.
- Download all dependencies
```bash
 go mod download
```

- Start postgresql server
```bash
make postgres PG_CONTAINER_NAME=<name>
```
- Create microbank database
```bash
make createdb
```
- Start redis server 
```bash
make redis 
```

- Start server
```bash
make server
```
* Query the endpoints using the **microbank.rest.http** file or any rest client out there

#### Run gRPC Gateway server
- Ensure this line of code is commented out in the **main** function
```go
runGinServer(config, store)
```
- Ensure these lines of code are uncommented in the **main** function
```go
go runGatewayServer(config, store)
runGrpcServer(config, store)
```
- Comment out the environment variables for docker environment and uncomment the ones for localhost in app.env
- Download all dependencies
```bash
go mod download
```

- Start postgresql server
```bash
make postgres PG_CONTAINER_NAME=<name>
```

- Create microbank database
```bash
make createdb
```

- Start redis server 
```bash
make redis 
```

- Start server
```bash
make server
```
* Query the endpoints using the **microbank.gateway.http** file or any rest client out there


#### Run only gRPC server
- Ensure these lines of code are commented out in the **main** function
```go
runGinServer(config, store)
go runGatewayServer(config, store)
```
- Ensure this line of code is uncommented in the **main** function
```go
runGrpcServer(config, store)
```
- Download the [evans gRPC client](https://intelops.ai/blog/evans-cli-a-go-grpc-client/#installation-of-evans-cli) on your machine
- Run the command
```bash
make evans
```