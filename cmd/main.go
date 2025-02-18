package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	app "github.com/Tairascii/google-docs-organization/internal"
	"github.com/Tairascii/google-docs-organization/internal/app/handler"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org/repo"
	"github.com/Tairascii/google-docs-organization/internal/app/service/user"
	"github.com/Tairascii/google-docs-organization/internal/app/usecase"
	"github.com/Tairascii/google-docs-organization/internal/db"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := app.LoadConfigs()
	dbSettings := db.Settings{
		Host:          cfg.Repo.Host,
		Port:          cfg.Repo.Port,
		User:          cfg.Repo.User,
		Password:      cfg.Repo.Password,
		DbName:        cfg.Repo.DBName,
		Schema:        cfg.Repo.Schema,
		AppName:       cfg.Repo.AppName,
		MaxIdleConns:  cfg.Repo.MaxIdleConns,
		MaxOpenConns:  cfg.Repo.MaxOpenConns,
		MigrateSchema: true,
	}

	sqlxDb, err := db.Connect(dbSettings)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
		return
	}
	defer func(sqlxDb *sqlx.DB) {
		if err := sqlxDb.Close(); err != nil {
			log.Fatalf("failed to close connection to db: %v", err)
		}
	}(sqlxDb)

	grpcAddr := fmt.Sprintf("%s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
	grpcClient, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcClient.Close()

	orgRepo := repo.New(sqlxDb)
	orgSrv := org.New(orgRepo)

	usrSrv := user.NewUserService(grpcClient)

	orgUC := usecase.NewOrgUseCase(orgSrv, usrSrv)

	useCase := app.UseCase{Org: orgUC}
	DI := &app.DI{UseCase: useCase}
	handlers := handler.NewHandler(DI)

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Something went wrong while runing server %s", err.Error())
		}
	}()

	log.Printf("Listening on port: %s", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Something went wrong while shutting down server %s", err.Error())
	}
}
