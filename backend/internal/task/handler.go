package task

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-assessment/backend/internal/httputil"
	"task-assessment/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *TaskService
}

func NewHandler(service *TaskService) *Handler {
	return &Handler{service: service}
}

type taskRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required"`
	Assignee    string `json:"assignee"`
}

func (h *Handler) List(c *gin.Context) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	limit := parsePositiveInt(c.DefaultQuery("limit", "10"), 10)
	if limit > 100 {
		limit = 100
	}

	status := strings.TrimSpace(c.Query("status"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	assignee := strings.TrimSpace(c.Query("assignee"))
	sort := strings.TrimSpace(c.Query("sort"))

	if status != "" && !isValidStatus(status) {
		httpError(c, http.StatusBadRequest, "INVALID_STATUS", "status must be todo, in_progress, or done")
		return
	}

	items, cacheHit, err := h.service.List(c.Request.Context(), models.ListParams{
		Page:     page,
		Limit:    limit,
		Status:   status,
		Keyword:  keyword,
		Assignee: assignee,
		Sort:     sort,
	})
	if err != nil {
		httpError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to load tasks")
		return
	}

	c.Header("X-Cache", map[bool]string{true: "HIT", false: "MISS"}[cacheHit])
	c.JSON(http.StatusOK, items)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpError(c, http.StatusNotFound, "TASK_NOT_FOUND", "task not found")
		return
	}
	if err != nil {
		httpError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to load task")
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Create(c *gin.Context) {
	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpError(c, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	item, err := h.service.Create(c.Request.Context(), models.CreateTaskInput(req))
	switch {
	case errors.Is(err, ErrDuplicateTitle):
		httpError(c, http.StatusConflict, "DUPLICATE_TITLE", "a task with this title already exists")
	case err != nil:
		httpError(c, http.StatusBadRequest, "INVALID_TASK", err.Error())
	default:
		c.JSON(http.StatusCreated, item)
	}
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpError(c, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, models.UpdateTaskInput(req))
	switch {
	case errors.Is(err, ErrNotFound):
		httpError(c, http.StatusNotFound, "TASK_NOT_FOUND", "task not found")
	case errors.Is(err, ErrDuplicateTitle):
		httpError(c, http.StatusConflict, "DUPLICATE_TITLE", "a task with this title already exists")
	case err != nil:
		httpError(c, http.StatusBadRequest, "INVALID_TASK", err.Error())
	default:
		c.JSON(http.StatusOK, item)
	}
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpError(c, http.StatusNotFound, "TASK_NOT_FOUND", "task not found")
			return
		}
		httpError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete task")
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		httpError(c, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func isValidStatus(status string) bool {
	switch status {
	case "todo", "in_progress", "done":
		return true
	default:
		return false
	}
}

func httpError(c *gin.Context, status int, code, message string) {
	httputil.JSONError(c, status, code, message)
}
