package handlers

import (
	"api-pos/model"
	"api-pos/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		h.GetAllCategories(w, r)
	case "POST":
		h.CreateCategory(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func (h *CategoryHandler) CategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Membersihkan path untuk mendapatkan ID
	path := strings.TrimPrefix(r.URL.Path, "/category/")
	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	switch r.Method {
	case "GET":
		h.GetCategoryByID(w, r, id)
	case "PUT":
		h.UpdateCategory(w, r, id)
	case "DELETE":
		h.DeleteCategory(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    categories,
	})
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Category not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    category,
	})
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	category, err := h.service.CreateCategory(newCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var updateData model.Category
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	category, err := h.service.UpdateCategory(id, updateData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteCategory(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category deleted successfully",
	})
}
