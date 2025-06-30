package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/muhadkprsnl/go-backend/internal/repositories"
)

// "local/qa-report/internal/repositories"
// ReleaseController handles release-related API endpoints
type ReleaseController struct {
	Repo *repositories.ReportRepository
}

// NewReleaseController initializes a new ReleaseController
func NewReleaseController(repo *repositories.ReportRepository) *ReleaseController {
	return &ReleaseController{Repo: repo}
}

// GetReleases handles GET /api/v1/releases requests with optional filters and pagination
func (c *ReleaseController) GetReleases(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status") // "on-time", "delayed", "all"
	skipStr := r.URL.Query().Get("skip")
	limitStr := r.URL.Query().Get("limit")

	// Convert skip and limit to int64
	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit == 0 {
		limit = 5 // Default limit
	}

	// Attach skip and limit to context
	ctx := context.WithValue(r.Context(), "skip", skip)
	ctx = context.WithValue(ctx, "limit", limit)

	// Fetch releases from repository
	releases, err := c.Repo.GetRecentFilteredReports(ctx, status)
	if err != nil {
		http.Error(w, "Error fetching releases: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(releases)
}
