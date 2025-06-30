package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"local/qa-report/internal/repositories"
)

type SummaryController struct {
	Repo *repositories.ReportRepository
}

func NewSummaryController(repo *repositories.ReportRepository) *SummaryController {
	return &SummaryController{Repo: repo}
}

func (c *SummaryController) GetSummary(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("startDate")
	endStr := r.URL.Query().Get("endDate")
	sprint := r.URL.Query().Get("sprint") // ✅ get sprint from query

	startDate, err1 := time.Parse("2006-01-02", startStr)
	endDate, err2 := time.Parse("2006-01-02", endStr)
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	summary, err := c.Repo.GetSummaryData(startDate, endDate, sprint) // ✅ pass it to repo
	if err != nil {
		http.Error(w, "Failed to fetch summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
