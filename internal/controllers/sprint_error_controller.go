package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

func (c *SummaryController) GetSprintErrorComparison(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("startDate")
	endStr := r.URL.Query().Get("endDate")

	startDate, err1 := time.Parse("2006-01-02", startStr)
	endDate, err2 := time.Parse("2006-01-02", endStr)
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	data, err := c.Repo.GetSprintErrorRates(startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to fetch error rates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
