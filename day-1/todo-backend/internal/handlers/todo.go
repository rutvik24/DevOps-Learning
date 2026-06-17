package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rutvik/todo-backend/internal/models"
	"gorm.io/gorm"
)

type TodoHandler struct {
	DB *gorm.DB
}

type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	if err := h.DB.Order("created_at DESC").Find(&todos).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch todos")
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch todo")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	if err := h.DB.Create(&todo).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create todo")
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	var todo models.Todo
	if err := h.DB.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch todo")
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != nil {
		if *req.Title == "" {
			writeError(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := h.DB.Save(&todo).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update todo")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	result := h.DB.Delete(&models.Todo{}, id)
	if result.Error != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete todo")
		return
	}
	if result.RowsAffected == 0 {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) Health(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.DB.DB()
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	if err := sqlDB.Ping(); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func parseID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	return uint(id), err
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
