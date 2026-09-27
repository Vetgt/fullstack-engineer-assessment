package main

import (
	"context"
	"fmt"
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

	tasks := []struct {
		title, description, status, assignee string
	}{
		{"Fix login validation", "Improve validation for login requests", "todo", "Betran"},
		{"Add dashboard pagination", "Add pagination to the dashboard task list", "in_progress", "Sarah"},
		{"Write API tests", "Cover update and search paths", "done", "Betran"},
		{"Improve Redis cache", "Review cache invalidation behavior", "todo", "Andi"},
		{"Polish task modal", "Improve task edit modal UX", "in_progress", "Sarah"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, task := range tasks {
		_, err := db.ExecContext(ctx, `
			INSERT INTO tasks (title, description, status, assignee)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE updated_at = updated_at`,
			task.title, task.description, task.status, task.assignee,
		)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Seed data inserted.")
}
