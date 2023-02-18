package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/Rexkizzy22/micro-bank/api"
	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/docs"
	"github.com/Rexkizzy22/micro-bank/gapi"
	_ "github.com/Rexkizzy22/micro-bank/gapi/statik"
	"github.com/Rexkizzy22/micro-bank/pb"
	"github.com/Rexkizzy22/micro-bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

// @securitydefinitions.apiKey ApiAuthKey
// @in                         header
// @name                       Authorization
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	dbConnString := config.FetchDBSource()
	dbDriver := config.FetchDBDriver()

	conn, err := sql.Open(dbDriver, dbConnString)
	if err != nil {
		log.Fatal("unable to connect to database: ", err)
	}

	// run db migration
	runMigration(config.MigrationURL, dbConnString)

	store := db.NewStore(conn)

	// Run HTTP server
	runGinServer(config, store)

	// Run GRPC Gateway & GRPC servers
	// go runGatewayServer(config, store)
	// runGrpcServer(config, store)
}

func runMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create migration instance: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run UP migration: ", err)
	}

	log.Println("db migrated successfully")
}

// RUN HTTP SERVER
func runGinServer(config util.Config, store db.Store) {
	setSwagger(config)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTP_ServerAddress)
	if err != nil {
		log.Fatal("unable to start server: ", err)
	}
}

// programmatically setting general swagger info
func setSwagger(config util.Config) {
	docs.SwaggerInfo.Title = "Micro Bank Rest API"
	docs.SwaggerInfo.Description = "A production-grade Go API that provides money transfer services between accounts of registered users"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = config.HTTP_ServerAddress
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMicroBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPC_ServerAddress)
	if err != nil {
		log.Fatal("unable to listen on grpc server: ", err)
	}

	log.Printf("grpc server listening at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("unable to start grpc server: ", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	// option to format grpc response jsonData to have snake-cased fields
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterMicroBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	// http network request handler
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// serves static swagger-ui assets from the grpc static asset folder
	// fs := http.FileServer(http.Dir("/gapi/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("failed to create statik asset server: ", err)
	}

	// serves static swagger assets from the statik server
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTP_ServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("unable to start HTTP gateway server: ", err)
	}
}
