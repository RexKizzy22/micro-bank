# Notes

The object of this project is to design, develop and deploy a production grade backend application, Microbank.

Microbank is an API that provides endpoints that enables the client to do the following operations:

1. Create and manage account - Owner, balance, currency
2. Record all balance changes - Create an account entry for each change
3. Facilitate money transfer transaction - Perform money transfer between two accounts consistently within a transaction

Below are the tools, key features and considerations for this application.

## DATABASE DESIGN

The first step is designing a SQL database schema using [dbdiagram.io](https://dbdiagram.io/home).

DB Diagram is a tool that enables us to save and share the database schema as a PDF/PNG diagram with others in the team.


## GENERATING CRUD CODE

We use code generation tools to generate code to run CRUD operations that interact 
with the database (serves as a repository pattern), we can create a schema in a target
database engine. In this project, we use [SQLC](https://sqlc.dev/).

Below shows how [SQLC](https://sqlc.dev/) compares with its alternatives in the golang ecosystem.

### [The Standard Database/SQL library](https://pkg.go.dev/database/sql)

Pros:

- Very fast and straightforward

Cons:

- Manual mapping SQL fields to variables
- Easy to make mistakes which are not caught until runtime.

### [GORM Library](gorm.io)

Pros:

- CRUD operations already implemented, thus, very short production code

Cons:

- Must learn to write queries using GORM's functions
- Runs slowly on high load

### [SQLX](https://github.com/jmoiron/sqlx)

GORM is quite opinionated so a more flexible option is SQLX.

Pros:

- Quite fast and easy to use
- Field mapping via query text and struct tags

Cons:

- Failure won't occur until runtime

### [SQLC](https://sqlc.dev/)

Pros:

- Uses the database/sql library under the hood
- Very fast and easy to use
- Uses the database/sql standard library
- Automatic code generation from manually written sql queries
- Evaluates the sql queries before generating code therefore errors are caught at compile time
- Generates `Querier` structs that makes it easy to mock the various database operations during testing

Cons:

- Only fully supports Postgres. MySQL is still experimental
  

## DATABASE MIGRATION

Over the course of the business lifetime, there will be changes made to the 
business model due to changes in business requirements. These changes will be 
reflected in the database schema and we will need to adapt our business model 
in a consistent and effective manner.

Database migration is a technique used to ensure that the schema is always 
adapted to the new business requirement in a consistent and reliable manner.

[The golang-migrate library](https://github.com/golang-migrate/migrate) is 
used in this project to run database migrations.


## API DESIGN

Microbank implements three kinds of application server. HTTP, gRPC and gRPC-Gateway servers.

The most noteworthy component of the API design is the ```server``` struct.

The HTTP server struct has the following properties:
```go
type Server struct {
   token token.Maker, // an interface that contains VerifyToken and CreateToken methods
   config *util.Config, // a struct that contains environment variables
   router *gin.Engine, // GIN router engine
   store db.Store // an interface that contains all functions that interact with the database
}
```

The gRPC server struct has the following properties:
```go
type Server struct {
   pb.UnimplementedMicroBankServer // aids forward compatibility, allows clients to peep services even before their implementation
   token token.Maker, // an interface that contains VerifyToken and CreateToken methods
   config *util.Config, // a struct that contains environment variables
   store db.Store // an interface that contains all functions that interact with the database
}
```

All handlers are methods of the ```server``` struct by way of receiver functions
```go
func (server *Server) HandlerName(ctx *gin.Context) {}
```

One of the most impressive things about this project is that the 
swagger assets are embedded in the the application server binary hence both 
the swagger docs and server are hosted and shipped together using the 
[**statik**](github.com/rakyll/statik/fs) library. This also means that 
we can use [**statik**](https://github.com/statik/statik) file server to 
serve client assets as well.

Microbank uses [**Viper**](https://github.com/viper), a robust configuration 
tool for managing configuration assets in local and remote environments

**Token-Based Authentication - Paseto over JWT**

- JSON Web Token **(JWT)** Authentication is not an entirely secure form of authentication
   - JWT gives developers too many hashing algorithms to choose from
   - Some algorithms are known to be vulnerable
     - RSA PKCSv1.5 padding oracle attack
     - ECDSA invalid-curve attack
   - A hacker can set the "alg" header to none
   - A hacker can set the "alg" header to "HS256" while the server normally verifies the token with a RSA public key

- Platform-Agnostic Security Token **(PASETO)** Authentication is a superior form of authentication due to the
following reasons:
  - No more "alg" header or "none" algorithm
  - Developers do not need to choose the hashing algorithm
  - Developers only need to choose the version of PASETO to use
  - Every version has one strong cipher suite
     - Only 2 most recent versions are accepted
     - v1 [compactable with legacy systems]
     - v2 [recommended]
  - Everything is authenticated
  - Encrypted payload for local use **(symmetric key)**

## UNIT TESTING

Golang was designed with unit testing in mind therefore making it part of the language binary. 
It has a robust unit testing convention where a ```TestMain``` function defined for each package 
in an application serves as an entry point for the tests in that package and the test files 
are colocated in the same directory as the files containing the code they are testing. 
The ```TestMain``` function usually takes in a parameter typically named ```m``` with type ```*testing.M``` 
which has a ```Run()``` method that should be called in an enclosing ```os.Exit()``` function call, 
thus ```os.Exit(m.Run())``` should be the last line of code in the ```TestMain``` function. 
This reports the success or failure of the tests ran in the package.

The [testify library](https://github.com/stretchr/testify) has several sub-packages containing
 testing utilities. We make use of the ```require``` sub-package which serves as a test-assertion 
 utility in this project. It is a mature and easy to use library.

## MOCKDB FOR UNIT TESTING

It is usually preferable to use a mock database to test all endpoints of an application. The benefits are:

- Independent tests - isolate test data to avoid conflicts especially when testing in a big project with a large codebase.
- Faster tests - reduce a lot of time talking to the database
- 100% coverage - enables us to write code that achieves 100% coverage and easily set up edge cases (unexpected errors)

Testing our API with a mock database is reliable and good enough mostly because our real db store is already tested.
All that is required is for the mock DB to implement the same interface as the real DB.

There are 2 ways to mock a database:
1. Use a fake DB (MEMORY) - implement a fake version of DB; store data in memory.
   - This is easy to implement but requires us to write more code which is time consuming for both development and maintenance.
2. Use DB stubs ([GOMOCK](https://github.com/golang/mock)) - Generate and build DB stubs that return hard-coded values

We use the db stubs in this project and write our unit tests in a table-driven manner 
to make the tests robust and more systematic


## TRANSACTIONS

Transactions are a single unit of database interaction that is made up of multiple operations.

To transfer $10 from account1 to account2, we have to:
1. Create a transfer record with amount = 10
2. Create an account entry for account1 with amount = -10
3. Create an account entry for account2 with amount = +10
4. Subtract $10 from the balance of account1
5. Add $10 to the balance of account2

**Why we need database transactions**
- To provide reliable and consistent unit of work, even in case of system failures
- To provide isolation between prograns that access the database concurrently

**Transactions need to satisfy the ACID properties in order to achieve its purpose:**
1. **Atomicity**: Either all the operations complete successfully or the transactions fail and the db is unchanged
2. **Consistency**: The db state must be valid after every transaction. All constraints must be satisfied
3. **Isolation**: Concurrent transactions must not affect each other
4. **Durabiltiy**: Data successfully written by a transaction must be recorded in a persistent storage

**Database transactions usually have 4 isolation levels as defined by the American National Standard Institute (ANSI). These levels define when a change made by one transaction is visible to the other**

- Read Uncommitted: Transactions can see data written by other uncommitted transactions
- Read Committed: Transactions can only see data written by other committed transactions
- Repeatable Read: Same read query always returns same results
- Serializable: Can achieve same result if execute transactions serially in some order instead of concurrently

**These levels are usually associated with 4 read phenomena:**

- **Dirty read**: A transaction reads data written by other concurrent uncommitted transaction
- **Non-repeatable read**: A transaction reads the same row twice and sees different values because it has been modified by other committed transaction
- **Phantom read**: A transation re-executes a query to find rows that satisfy a condition and sees a different set of rows, due to changes by other committed transaction
- **Serialization anomal**: The result of a group of concurrent committed transactions is impossible to achieve if we run them sequentially in any order without overlapping

**MySQL vs PostgresSQL**
MySQL
- MYSQL uses locking mechanism to prevent serialization anomaly
- By default, it is Repeatable Read

PostgreSQL
- PostgreSQL uses read/write dependency detection to prevent serialization anomaly
- By default, it is Read Committed

**Note:**
1. Implement retry mechanisms when using high isolation levels as there could be errors, timeouts or deadlocks
2. Each database might implement isolation level differently - Read documentation


## GitHub Actions

**Workflow**
- A Workflow is an automated procedure made up of one or more jobs. They are triggered by events, scheduled or manually triggered.
- Runners are servers for running jobs. They run jobs one at a time and can be hosted on GitHub or self-hosted.
- A Runner reports logs, progress and results to GitHub.

**Job**
- A Job is a set of steps that are executed on the same runner
- Normal jobs run in parallel
- Dependent jobs run serially

**Step**
- A Step is an individual task run serially within a Job
- A Step contains one or more actions

**Action**
- An Action is a standalone command that runs serially within a Step
- An Action can be reused
