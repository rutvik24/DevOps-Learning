package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rutvik/todo-backend/internal/config"
	"github.com/rutvik/todo-backend/internal/database"
	"github.com/rutvik/todo-backend/internal/handlers"
	"github.com/rutvik/todo-backend/internal/middleware"
	"github.com/rutvik/todo-backend/internal/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if err := migrations.Up(db); err != nil {
		log.Fatalf("migrations up: %v", err)
	}

	todoHandler := &handlers.TodoHandler{DB: db}

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", todoHandler.Health).Methods(http.MethodGet)
	api.HandleFunc("/todos", todoHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/todos", todoHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/todos/{id}", todoHandler.Get).Methods(http.MethodGet)
	api.HandleFunc("/todos/{id}", todoHandler.Update).Methods(http.MethodPut, http.MethodPatch)
	api.HandleFunc("/todos/{id}", todoHandler.Delete).Methods(http.MethodDelete)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      middleware.CORS(cfg.CORSOrigin)(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
