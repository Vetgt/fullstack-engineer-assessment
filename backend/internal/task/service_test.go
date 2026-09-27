package task

import (
	"context"
	"testing"
	"time"

	"task-assessment/backend/internal/cache"
	"task-assessment/backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// dbRepoAdapter keeps the test focused on service -> repository -> SQL behavior.
// The production repository is used below instead of a separate fake.
func setupTestRepo(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return NewRepository(db), mock, func() {
		_ = db.Close()
	}
}

func TestUpdateTask(t *testing.T) {
	repo, mock, cleanup := setupTestRepo(t)
	defer cleanup()

	mock.ExpectExec("UPDATE tasks").
		WithArgs(
			"Updated title",
			"Updated description",
			"in_progress",
			"Betran",
			uint64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT id, title").
		WithArgs(uint64(1)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"title",
				"COALESCE(description, '')",
				"status",
				"COALESCE(assignee, '')",
				"created_at",
				"updated_at",
			}).AddRow(
				1,
				"Updated title",
				"Updated description",
				"in_progress",
				"Betran",
				time.Date(2026, 9, 27, 10, 0, 0, 0, time.Local),
				time.Date(2026, 9, 27, 10, 1, 0, 0, time.Local),
			),
		)

	item, err := repo.Update(
		context.Background(),
		1,
		models.UpdateTaskInput{
			Title:       "Updated title",
			Description: "Updated description",
			Status:      "in_progress",
			Assignee:    "Betran",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "Updated title", item.Title)
	require.Equal(t, "in_progress", item.Status)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchQueryIncludesKeyword(t *testing.T) {
	repo, mock, cleanup := setupTestRepo(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM tasks`).
		WithArgs("%login%", "%login%").
		WillReturnRows(
			sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1),
		)

	mock.ExpectQuery("SELECT id, title").
		WithArgs("%login%", "%login%", 10, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"title",
				"COALESCE(description, '')",
				"status",
				"COALESCE(assignee, '')",
				"created_at",
				"updated_at",
			}).AddRow(
				1,
				"Fix login",
				"Login issue",
				"todo",
				"Betran",
				time.Date(2026, 9, 27, 10, 0, 0, 0, time.Local),
				time.Date(2026, 9, 27, 10, 1, 0, 0, time.Local),
			),
		)

	result, err := repo.List(
		context.Background(),
		models.ListParams{
			Keyword: "login",
			Page:    1,
			Limit:   10,
		},
	)

	require.NoError(t, err)
	require.Len(t, result.Items, 1)
	require.Equal(t, "Fix login", result.Items[0].Title)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheInvalidationAfterMutation(t *testing.T) {
	mini := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mini.Addr(),
	})
	defer rdb.Close()

	c := cache.NewRedisCache(
		mini.Addr(),
		"",
		0,
		60_000_000_000,
	)

	// Use the production cache client because its invalidation behavior
	// is what the assessment asks for.
	_ = rdb

	ctx := context.Background()
	key := cache.ListCacheKey("page=1&limit=10")

	require.NoError(
		t,
		c.SetTaskList(
			ctx,
			key,
			models.TaskList{
				Page:  1,
				Limit: 10,
				Total: 1,
			},
		),
	)

	require.True(t, mini.Exists(key))

	require.NoError(t, c.InvalidateTaskLists(ctx))
	require.False(t, mini.Exists(key))
}
