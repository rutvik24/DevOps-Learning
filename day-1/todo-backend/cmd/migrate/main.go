package main

import (
	"flag"
	"log"

	"github.com/rutvik/todo-backend/internal/config"
	"github.com/rutvik/todo-backend/internal/database"
	"github.com/rutvik/todo-backend/internal/migrations"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up, down, or down-all")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	switch *direction {
	case "up":
		if err := migrations.Up(db); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied successfully")
	case "down":
		if err := migrations.Down(db); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("last migration rolled back successfully")
	case "down-all":
		if err := migrations.DownAll(db); err != nil {
			log.Fatalf("migrate down-all: %v", err)
		}
		log.Println("all migrations rolled back successfully")
	default:
		log.Fatalf("unknown direction: %s (use up, down, or down-all)", *direction)
	}
}
