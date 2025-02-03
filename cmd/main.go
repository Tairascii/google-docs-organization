package main

import (
	"context"
	app "github.com/Tairascii/google-docs-organization/internal"
	"github.com/Tairascii/google-docs-organization/internal/app/handler"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org/repo"
	"github.com/Tairascii/google-docs-organization/internal/app/usecase"
	"github.com/Tairascii/google-docs-organization/internal/db"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//TODO add envs
	dbSettings := db.Settings{
		Host:          "localhost",
		Port:          "5432",
		User:          "admin",
		Password:      "12345",
		DbName:        "google_doc_organization",
		Schema:        "google_doc_organization_schema",
		AppName:       "google_doc_organization",
		MaxIdleConns:  2,
		MaxOpenConns:  5,
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

	orgRepo := repo.New(sqlxDb)
	orgSrv := org.New(orgRepo)

	orgUC := usecase.NewOrgUseCase(orgSrv)

	useCase := app.UseCase{Org: orgUC}
	DI := &app.DI{UseCase: useCase}
	handlers := handler.NewHandler(DI)

	srv := &http.Server{
		Addr:         ":8000", // TODO add .env
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Something went wrong while runing server %s", err.Error())
		}
	}()

	log.Println("Listening on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Something went wrong while shutting down server %s", err.Error())
	}
}
