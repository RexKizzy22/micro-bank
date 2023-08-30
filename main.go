package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Rexkizzy22/micro-bank/api"
	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/gapi"
	_ "github.com/Rexkizzy22/micro-bank/gapi/statik"
	"github.com/Rexkizzy22/micro-bank/mail"
	"github.com/Rexkizzy22/micro-bank/pb"
	docs "github.com/Rexkizzy22/micro-bank/swagger"
	"github.com/Rexkizzy22/micro-bank/task"
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

// @securitydefinitions.apiKey  ApiAuthKey
// @in                          header
// @name                        Authorization
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("unable to load config: %s", err)
	}

	if config.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	dbConnString := config.FetchDBSource()
	dbDriver := config.FetchDBDriver()

	conn, err := sql.Open(dbDriver, dbConnString)
	if err != nil {
		log.Fatal().Msgf("unable to connect to database: %s", err)
	}

	// run db migration
	runMigration(config.MigrationURL, dbConnString)

	store := db.NewStore(conn)

	// setup redis queue for server integration
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := task.NewRedisTaskDistributor(redisOpt)

	// run redis task processor
	go runTaskProcessor(config, &redisOpt, store)

	// Run HTTP server
	// runGinServer(config, store)

	// Run gRPC-Gateway & gRPC servers
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msgf("cannot create migration instance: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("failed to run UP migration: %s", err)
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt *asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := task.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Err(err).Msg("failed to start task processor")
	}
}

// RUN HTTP SERVER
func runGinServer(config util.Config, store db.Store) {
	setSwagger(config)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	err = server.Start(config.HTTP_ServerAddress)
	if err != nil {
		log.Fatal().Msgf("unable to start server: %s", err)
	}
}

// programmatically setting general swagger info
func setSwagger(config util.Config) {
	docs.SwaggerInfo_swagger.Title = "Micro Bank Rest API"
	docs.SwaggerInfo_swagger.Description = "A production-grade Go API that provides money transfer services between accounts of registered users"
	docs.SwaggerInfo_swagger.Version = "1.0.0"
	docs.SwaggerInfo_swagger.Host = config.HTTP_ServerAddress
	docs.SwaggerInfo_swagger.BasePath = "/"
	docs.SwaggerInfo_swagger.Schemes = []string{"http"}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor task.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterMicroBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPC_ServerAddress)
	if err != nil {
		log.Fatal().Msgf("unable to listen on grpc server: %s", err)
	}

	log.Info().Msgf("grpc server listening at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("unable to start grpc server: %s", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor task.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	// option to format gRPC response jsonData to have snake-cased fields
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
		log.Fatal().Msgf("cannot register handler server: %s", err)
	}

	// http network request handler
	mux := http.NewServeMux()
	handler := gapi.HTTPLogger(mux)
	mux.Handle("/", handler)

	// serves static swagger-ui assets from the gRPC static asset folder
	// fs := http.FileServer(http.Dir("/gapi/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msgf("failed to create statik asset server: %s", err)
	}

	// serves static swagger assets from the statik server
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTP_ServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener: %s", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Msgf("unable to start HTTP gateway server: %s", err)
	}
}
