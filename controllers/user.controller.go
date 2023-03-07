package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/stockk-go/models"
	"github.com/hiennguyen9874/stockk-go/repository"
	"gorm.io/gorm"
)

type UserHandler struct {
	db   *gorm.DB
	repo repository.UserRepo
}

func NewUserHandler(db *gorm.DB, repo repository.UserRepo) *UserHandler {
	return &UserHandler{
		db:   db,
		repo: repo,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	user, err := h.repo.Get(h.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(h.repo.GetModelResponse(user))
	render.Respond(w, r, h.repo.GetModelResponse(user))
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	users, err := h.repo.GetMulti(h.db, limit, offset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(&users)
	render.Respond(w, r, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.UserCreate)

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.repo.Create(h.db, h.repo.GetModel(user))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(h.repo.GetModelResponse(newUser))
	render.Respond(w, r, h.repo.GetModelResponse(newUser))
}
