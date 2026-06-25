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
	authHandler "github.com/rusilkoirala/student-api/internal/http/handlers/auth"
	classHandler "github.com/rusilkoirala/student-api/internal/http/handlers/class"
	studentHandler "github.com/rusilkoirala/student-api/internal/http/handlers/student"
	teacherHandler "github.com/rusilkoirala/student-api/internal/http/handlers/teacher"
	"github.com/rusilkoirala/student-api/internal/http/middleware"
	"github.com/rusilkoirala/student-api/internal/storage/sqlite"
	"github.com/rusilkoirala/student-api/internal/utils/response"
)

func main() {
	cfg := config.MustLoad()

	store, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env))

	// ── Public routes (no auth) ───────────────────────────────────────────────
	public := http.NewServeMux()
	public.HandleFunc("POST /api/auth/register", authHandler.Register(store, cfg.JWTSecret))
	public.HandleFunc("POST /api/auth/login", authHandler.Login(store, cfg.JWTSecret))

	// ── Protected routes (JWT required) ──────────────────────────────────────
	protected := http.NewServeMux()

	protected.HandleFunc("POST /api/students", studentHandler.Create(store))
	protected.HandleFunc("GET /api/students/{id}", studentHandler.GetById(store))
	protected.HandleFunc("GET /api/students", studentHandler.GetList(store))
	protected.HandleFunc("DELETE /api/students/{id}", studentHandler.Delete(store))

	protected.HandleFunc("POST /api/teachers", teacherHandler.Create(store))
	protected.HandleFunc("GET /api/teachers/{id}", teacherHandler.GetById(store))
	protected.HandleFunc("GET /api/teachers", teacherHandler.GetList(store))
	protected.HandleFunc("DELETE /api/teachers/{id}", teacherHandler.Delete(store))

	protected.HandleFunc("POST /api/classes", classHandler.Create(store))
	protected.HandleFunc("GET /api/classes/{id}", classHandler.GetById(store))
	protected.HandleFunc("GET /api/classes", classHandler.GetList(store))
	protected.HandleFunc("DELETE /api/classes/{id}", classHandler.Delete(store))

	protected.HandleFunc("GET /api/stats", func(w http.ResponseWriter, r *http.Request) {
		schoolId, ok := middleware.SchoolIDFromCtx(w, r)
		if !ok {
			return
		}
		stats, err := store.GetStats(schoolId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, stats)
	})

	// ── Root handler: CORS + route dispatch ───────────────────────────────────
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Try public first, then protected (with auth middleware)
		if r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
			public.ServeHTTP(w, r)
			return
		}

		middleware.RequireAuth(cfg.JWTSecret, protected).ServeHTTP(w, r)
	})

	server := http.Server{
		Addr:    cfg.Address,
		Handler: root,
	}

	fmt.Println("Server running on", cfg.Address)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	slog.Info("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", slog.String("error", err.Error()))
	}
	slog.Info("Server stopped")
}
