package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"task-assessment/backend/internal/config"
	"task-assessment/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.OpenMySQL(cfg.MySQLDSN())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, string(sqlBytes)); err != nil {
		panic(err)
	}

	fmt.Println("Database migration completed.")
}
