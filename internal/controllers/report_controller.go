package controllers

import (
	"encoding/json"
	"local/qa-report/internal/models"
	"local/qa-report/internal/repositories"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ReportController struct {
	repo   *repositories.ReportRepository
	logger *zap.Logger
}

func NewReportController(repo *repositories.ReportRepository, logger *zap.Logger) *ReportController {
	return &ReportController{
		repo:   repo,
		logger: logger,
	}
}

func (c *ReportController) CreateReport(w http.ResponseWriter, r *http.Request) {
	var report models.FormData
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		c.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, err := c.repo.CreateReport(report)
	if err != nil {
		c.logger.Error("Failed to create report", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Report created successfully",
		"id":      id,
	})
}

func (c *ReportController) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := c.repo.GetAllReports()
	if err != nil {
		c.logger.Error("Failed to get reports", zap.Error(err))
		http.Error(w, "Failed to fetch reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reports); err != nil {
		c.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to process data", http.StatusInternalServerError)
	}
}

func (c *ReportController) UpdateReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var report models.FormData
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		c.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := c.repo.UpdateReport(id, report); err != nil {
		c.logger.Error("Failed to update report", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Report updated successfully"))
}

func (c *ReportController) DeleteReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repo.DeleteReport(id); err != nil {
		c.logger.Error("Failed to delete report", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Report deleted successfully"))
}
