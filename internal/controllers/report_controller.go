package controllers

import (
	"encoding/json"
	"time"

	// "local/qa-report/internal/models"

	// "local/qa-report/internal/repositories"
	"net/http"

	"github.com/muhadkprsnl/go-backend/internal/models"
	"github.com/muhadkprsnl/go-backend/internal/repositories"

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

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"message": "Report created successfully",
	// 	"id":      id,
	// })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Report updated successfully",
		"id":      id, // optional — only if you want to return it
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

	var requestData struct {
		Sprint         string    `json:"sprint"`
		Version        string    `json:"version"`
		DueDate        time.Time `json:"dueDate"`
		CloseDate      time.Time `json:"closeDate"`
		TotalTestCases int       `json:"totalTestCases"`
		TotalBugs      int       `json:"totalBugs"`
		Developers     []struct {
			Name   string `json:"name"`
			Passed int    `json:"passed"`
			Failed int    `json:"failed"`
		} `json:"developers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		c.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Convert to the database model structure
	report := models.FormData{
		Sprint:        requestData.Sprint,
		Version:       requestData.Version,
		DueDate:       requestData.DueDate,
		CloseDate:     requestData.CloseDate,
		Totaltestcase: requestData.TotalTestCases,
		Totalbugs:     requestData.TotalBugs,
	}

	// Map developers if they exist
	if len(requestData.Developers) > 0 {
		report.Developer1 = requestData.Developers[0].Name
		report.D1Passed = requestData.Developers[0].Passed
		report.D1Failed = requestData.Developers[0].Failed
	}
	if len(requestData.Developers) > 1 {
		report.Developer2 = requestData.Developers[1].Name
		report.D2Passed = requestData.Developers[1].Passed
		report.D2Failed = requestData.Developers[1].Failed
	}

	if err := c.repo.UpdateReport(id, report); err != nil {
		c.logger.Error("Failed to update report", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the updated report
	updatedReport, err := c.repo.GetReportByID(id)
	if err != nil {
		c.logger.Error("Failed to fetch updated report", zap.Error(err))
		http.Error(w, "Failed to fetch updated data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedReport)
}

// func (c *ReportController) UpdateReport(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	var report models.FormData
// 	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
// 		c.logger.Error("Failed to decode request body", zap.Error(err))
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	if err := c.repo.UpdateReport(id, report); err != nil {
// 		c.logger.Error("Failed to update report", zap.Error(err))
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Report updated successfully"))
// }

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
