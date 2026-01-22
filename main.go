package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Kategori struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var ProdukList = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 5000, Stok: 20},
	{ID: 2, Nama: "Teh Pucuk", Harga: 3000, Stok: 42},
}

var Category = []Kategori{
	{ID: 1, Name: "Makanan", Description: "Segala jenis makanan"},
	{ID: 2, Name: "Minuman", Description: "Minuman Kemasan, Teh, Kopi"},
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkHandlerByID)
	http.HandleFunc("/categories", categoryHandler)
	http.HandleFunc("/categories/", categoryHandlerByID)

	fmt.Println("Server running on at http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Gagal menjalankan server: %v\n", err)
	}
}

// Root Handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "API is running well",
	})
}

// Produk handler : Get semua produk dan Create Produk
func produkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		getProduk(w, r)
	case "POST":
		createProduk(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Produk Handler By Id, Get Detail, Edit Produk, Delete Produk
func produkHandlerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extract ID dari url
	path := strings.TrimPrefix(r.URL.Path, "/produk/")
	id, err := strconv.Atoi(path)

	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getProdukDetail(w, r, id)
	case "PUT":
		editProduk(w, r, id)
	case "DELETE":
		deleteProduk(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handler untuk get semua produk
func getProduk(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    ProdukList,
	})
}

// handler untuk get detail
func getProdukDetail(w http.ResponseWriter, r *http.Request, id int) {
	for _, p := range ProdukList {
		if p.ID == id {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    p,
			})
			return
		}
	}
	http.Error(w, "Produk Not Found", http.StatusNotFound)

}

// handler POST - Create Produk
func createProduk(w http.ResponseWriter, r *http.Request) {
	var newProduk Produk
	err := json.NewDecoder(r.Body).Decode(&newProduk)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Generate ID by length produk ++
	newProduk.ID = len(ProdukList) + 1
	ProdukList = append(ProdukList, newProduk)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    newProduk,
		"message": "Produk Created Successfully",
	})
}

// handler PUT - Edit Produk
func editProduk(w http.ResponseWriter, r *http.Request, id int) {
	// get data dari request
	var updatedProduk Produk
	err := json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// loop produk, cari id
	for i, p := range ProdukList {
		if p.ID == id {
			ProdukList[i].Nama = updatedProduk.Nama
			ProdukList[i].Harga = updatedProduk.Harga
			ProdukList[i].Stok = updatedProduk.Stok
			ProdukList[i].ID = id // keep origin id

			json.NewEncoder(w).Encode(map[string]interface{}{
				"succes":  true,
				"data":    ProdukList[i],
				"message": "Produk updated successfully",
			})
			return
		}
	}
	http.Error(w, "Produk Not Found", http.StatusNotFound)

}

// handler Delete -- Delete Produk
func deleteProduk(w http.ResponseWriter, r *http.Request, id int) {
	// loop produk , cari id
	for i, p := range ProdukList {
		if p.ID == id {
			ProdukList = append(ProdukList[:i], ProdukList[i+1:]...)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Produk deleted successfully",
			})
			return
		}
	}
	http.Error(w, "Produk Not Found", http.StatusNotFound)
}

// ===================== Handler Categories ==================== //
// Category handler : Get semua Category dan Create Category
func categoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		getCategory(w, r)
	case "POST":
		createCategory(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Category Handler By Id, Get Detail, Edit Category, Delete Category
func categoryHandlerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extract ID dari url
	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(path)

	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getCategoryDetail(w, r, id)
	case "PUT":
		editCategory(w, r, id)
	case "DELETE":
		deleteCategory(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handler untuk get semua produk
func getCategory(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    Category,
	})
}

// handler untuk get detail
func getCategoryDetail(w http.ResponseWriter, r *http.Request, id int) {
	for _, p := range Category {
		if p.ID == id {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    p,
			})
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)
}

// handler POST - Create Category
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Kategori
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Generate ID by length category ++
	newCategory.ID = len(Category) + 1
	Category = append(Category, newCategory)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    newCategory,
		"message": "Category Created Successfully",
	})
}

// handler PUT - Edit Category
func editCategory(w http.ResponseWriter, r *http.Request, id int) {
	// get data dari request
	var updatedCategory Kategori
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// loop Category, cari id
	for i, p := range Category {
		if p.ID == id {
			Category[i].Name = updatedCategory.Name
			Category[i].Description = updatedCategory.Description
			Category[i].ID = id // keep origin id

			json.NewEncoder(w).Encode(map[string]interface{}{
				"succes":  true,
				"data":    Category[i],
				"message": "Category updated successfully",
			})
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)

}

// handler Delete -- Delete Category
func deleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	// loop category , cari id
	for i, p := range Category {
		if p.ID == id {
			Category = append(Category[:i], Category[i+1:]...)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Category deleted successfully",
			})
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)
}
