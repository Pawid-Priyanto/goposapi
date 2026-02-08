package handlers

import (
	"api-pos/model"
	"api-pos/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req model.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	now := time.Now().Format("2006-01-02")

	if startDate == "" {
		startDate = now
	}
	if endDate == "" {
		endDate = now
	}
	fmt.Printf("Querying for dates: %s to %s\n", startDate, endDate)

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {

		http.Error(w, "failed to generated "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content_Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
