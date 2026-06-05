package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rusilkoirala/student-api/internal/config"
	"github.com/rusilkoirala/student-api/internal/http/handlers/student"
	"github.com/rusilkoirala/student-api/internal/storage/sqlite"
)

func main() {

	// Load config
	cfg := config.MustLoad()
	// Db setup
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// Router setup

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.Create(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	router.HandleFunc("GET /api/students", student.GetList(storage))

	// Setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Println("Server is running")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}

	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
