package task

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"task-assessment/backend/internal/cache"
	"task-assessment/backend/internal/models"
)

type TaskRepository interface {
	List(context.Context, models.ListParams) (models.TaskList, error)
	GetByID(context.Context, uint64) (models.Task, error)
	Create(context.Context, models.CreateTaskInput) (models.Task, error)
	Update(context.Context, uint64, models.UpdateTaskInput) (models.Task, error)
	SoftDelete(context.Context, uint64) error
}

type TaskService struct {
	repo  TaskRepository
	cache *cache.RedisCache
}

func NewService(repo TaskRepository, cacheClient *cache.RedisCache) *TaskService {
	return &TaskService{repo: repo, cache: cacheClient}
}

func (s *TaskService) List(ctx context.Context, params models.ListParams) (models.TaskList, bool, error) {
	cacheKey := cache.ListCacheKey(buildListQuery(params))

	var cached models.TaskList
	if hit, err := s.cache.GetTaskList(ctx, cacheKey, &cached); err == nil && hit {
		return cached, true, nil
	}

	result, err := s.repo.List(ctx, params)
	if err != nil {
		return models.TaskList{}, false, err
	}

	if err := s.cache.SetTaskList(ctx, cacheKey, result); err != nil {
		// A cache failure should not break the API response.
		return result, false, nil
	}
	return result, false, nil
}

func (s *TaskService) Get(ctx context.Context, id uint64) (models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) Create(ctx context.Context, input models.CreateTaskInput) (models.Task, error) {
	input.Title = normalizeTitle(input.Title)
	if err := validateInput(input.Title, input.Status); err != nil {
		return models.Task{}, err
	}
	item, err := s.repo.Create(ctx, input)
	if err != nil {
		return models.Task{}, err
	}
	_ = s.cache.InvalidateTaskLists(ctx)
	return item, nil
}

func (s *TaskService) Update(ctx context.Context, id uint64, input models.UpdateTaskInput) (models.Task, error) {
	input.Title = normalizeTitle(input.Title)
	if err := validateInput(input.Title, input.Status); err != nil {
		return models.Task{}, err
	}
	item, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return models.Task{}, err
	}
	_ = s.cache.InvalidateTaskLists(ctx)
	return item, nil
}

func (s *TaskService) SoftDelete(ctx context.Context, id uint64) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.InvalidateTaskLists(ctx)
	return nil
}

func normalizeTitle(title string) string {
	return strings.TrimSpace(title)
}

func validateInput(title, status string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if len(title) > 255 {
		return fmt.Errorf("title must be 255 characters or fewer")
	}
	switch status {
	case "todo", "in_progress", "done":
		return nil
	default:
		return fmt.Errorf("status must be todo, in_progress, or done")
	}
}

func buildListQuery(params models.ListParams) string {
	q := url.Values{}
	q.Set("page", strconv.Itoa(params.Page))
	q.Set("limit", strconv.Itoa(params.Limit))
	q.Set("sort", params.Sort)
	if params.Keyword != "" {
		q.Set("keyword", params.Keyword)
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Assignee != "" {
		q.Set("assignee", params.Assignee)
	}
	return q.Encode()
}
