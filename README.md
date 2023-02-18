# Microbank

### A minimalist implementation of a banking server, serving the following purposes:

* creating a banking customer
* creating an account for a particular customer using one of the accepted currencies
* transfering money between customer accounts
* retrieving all accounts created by the customer in the current session
* retrieving one account created by the customer in the current session

**Read [NOTE.md](NOTE.md) for detailed information about this project**

## Run Microbank

**The easiest way to run Microbank server is to spin up its docker services**

* Download [Docker](https://www.docker.com/products/docker-desktop/)
* Start services

```bash
docker-compose up --build
```
**Another way is to start the server on localhost**

### Pre-requisites for running Microbank on a Macbook

1. [Homebrew](https://brew.sh/)
2. PostgreSQL. Run `brew install postgresql`
3. [Go](https://go.dev/dl/)
4. golang-migrate. Run `brew install golang-migrate`

#### Run HTTP server

* Start postgresql server

```bash
brew services start postgresql
```
* Provide environment variables as hinted by the [example.env](example.env) file
* Download all dependencies

```bash
go install
```

* Start server
```bash
make server
```
* Call REST API using HTTP clients like Postman

#### Run GRPC Gateway & GRPC servers
* Comment out line `48` in [main.go](main.go) file
* Uncomment lines `51` and `52` in [main.go](main.go) file
* (Re)Start server
```bash
make server
```