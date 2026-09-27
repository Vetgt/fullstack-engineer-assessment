package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"task-assessment/backend/internal/models"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, params models.ListParams) (models.TaskList, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]any, 0, 8)

	if params.Status != "" {
		where = append(where, "status = ?")
		args = append(args, params.Status)
	}
	if params.Keyword != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		keyword := "%" + params.Keyword + "%"
		args = append(args, keyword, keyword)
	}
	if params.Assignee != "" {
		where = append(where, "assignee LIKE ?")
		args = append(args, "%"+params.Assignee+"%")
	}

	orderBy, err := safeSort(params.Sort)
	if err != nil {
		return models.TaskList{}, err
	}

	countSQL := "SELECT COUNT(*) FROM tasks WHERE " + strings.Join(where, " AND ")
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return models.TaskList{}, err
	}

	offset := (params.Page - 1) * params.Limit
	querySQL := `
		SELECT id, title, COALESCE(description, ''), status, COALESCE(assignee, ''), created_at, updated_at
		FROM tasks
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY ` + orderBy + `
		LIMIT ? OFFSET ?`
	queryArgs := append(append([]any{}, args...), params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, querySQL, queryArgs...)
	if err != nil {
		return models.TaskList{}, err
	}
	defer rows.Close()

	items := make([]models.Task, 0, params.Limit)
	for rows.Next() {
		var item models.Task
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Status, &item.Assignee, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return models.TaskList{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return models.TaskList{}, err
	}

	return models.TaskList{
		Items:      items,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (models.Task, error) {
	const query = `
		SELECT id, title, COALESCE(description, ''), status, COALESCE(assignee, ''), created_at, updated_at
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL`

	var item models.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.Status,
		&item.Assignee,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Task{}, ErrNotFound
	}
	if err != nil {
		return models.Task{}, err
	}
	return item, nil
}

func (r *Repository) Create(ctx context.Context, input models.CreateTaskInput) (models.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status, assignee)
		VALUES (?, NULLIF(?, ''), ?, NULLIF(?, ''))`

	result, err := r.db.ExecContext(ctx, query, input.Title, input.Description, input.Status, input.Assignee)
	if err != nil {
		return models.Task{}, normalizeDBError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, err
	}
	return r.GetByID(ctx, uint64(id))
}

func (r *Repository) Update(ctx context.Context, id uint64, input models.UpdateTaskInput) (models.Task, error) {
	const query = `
		UPDATE tasks
		SET title = ?, description = NULLIF(?, ''), status = ?, assignee = NULLIF(?, '')
		WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, input.Title, input.Description, input.Status, input.Assignee, id)
	if err != nil {
		return models.Task{}, normalizeDBError(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return models.Task{}, err
	}
	if affected == 0 {
		return models.Task{}, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) SoftDelete(ctx context.Context, id uint64) error {
	const query = `UPDATE tasks SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func normalizeDBError(err error) error {
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrDuplicateTitle
	}
	return err
}

func safeSort(sort string) (string, error) {
	switch sort {
	case "created_at_asc":
		return "created_at ASC, id ASC", nil
	case "created_at_desc", "":
		return "created_at DESC, id DESC", nil
	case "title_asc":
		return "title ASC, id ASC", nil
	case "title_desc":
		return "title DESC, id DESC", nil
	case "status_asc":
		return "status ASC, id ASC", nil
	case "status_desc":
		return "status DESC, id DESC", nil
	default:
		return "", fmt.Errorf("unsupported sort: %s", sort)
	}
}
