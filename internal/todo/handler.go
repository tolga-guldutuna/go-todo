package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", h.handleTodos)
	mux.HandleFunc("/todos/", h.handleTodoByID) // /todos/1 gibi
}

func (h *Handler) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	// URL: /todos/123 → "123"ü al
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetOne(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetTodos godoc
// @Summary      List all todos
// @Tags         todos
// @Produce      json
// @Success      200  {array}   Todo
// @Router       /todos [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := h.svc.GetTodos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

// GetTodo godoc
// @Summary      Get todo by ID
// @Tags         todos
// @Produce      json
// @Param        id   path      int   true  "Todo ID"
// @Success      200  {object}  Todo
// @Router       /todos/{id} [get]
func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request, id int) {
	todo, err := h.svc.GetTodo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// CreateTodo godoc
// @Summary      Create a new todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body      createTodoRequest  true  "Todo to create"
// @Success      201   {object}  Todo
// @Router       /todos [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	todo, err := h.svc.AddTodo(r.Context(), req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

type updateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsDone      *bool   `json:"isDone,omitempty"`
}

// UpdateTodo godoc
// @Summary      Update a todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Todo ID"
// @Param        todo  body      updateTodoRequest  true  "Todo to update"
// @Success      200   {object}  Todo
// @Router       /todos/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.GetTodo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.IsDone != nil {
		todo.IsDone = *req.IsDone
	}

	updated, err := h.svc.UpdateTodo(
		r.Context(),
		todo.ID,
		todo.Title,
		todo.Description,
		todo.IsDone,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DeleteTodo godoc
// @Summary      Delete a todo
// @Tags         todos
// @Param        id   path      int  true  "Todo ID"
// @Success      204  "No Content"
// @Router       /todos/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.svc.DeleteTodo(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
