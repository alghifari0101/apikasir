package handlers

import (
	"apikasir/services"
	"encoding/json"
	"net/http"
)

type ReportHandler struct {
	Service *services.TransactionService
}

func NewReportHandler(service *services.TransactionService) *ReportHandler {
	return &ReportHandler{Service: service}
}

func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.Service.Repo.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetReportByRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.Service.Repo.GetReportByRange(r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
