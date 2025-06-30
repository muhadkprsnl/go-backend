package controllers

import (
	"encoding/json"
	"fmt"

	// "local/qa-report/internal/repositories"
	"net/http"
	"time"

	"github.com/muhadkprsnl/go-backend/internal/repositories"
	"go.uber.org/zap"
)

type ChartController struct {
	repo   *repositories.ChartRepository
	logger *zap.Logger
}

func NewChartController(repo *repositories.ChartRepository, logger *zap.Logger) *ChartController {
	return &ChartController{
		repo:   repo,
		logger: logger,
	}
}

// func (c *ChartController) GetDonutChart(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	query := r.URL.Query()

// 	sprint := query.Get("sprint")

// 	var startDatePtr, endDatePtr *time.Time
// 	layout := "2006-01-02"

// 	if startStr := query.Get("startDate"); startStr != "" {
// 		if startDate, err := time.Parse(layout, startStr); err == nil {
// 			startDatePtr = &startDate
// 		}
// 	}
// 	if endStr := query.Get("endDate"); endStr != "" {
// 		if endDate, err := time.Parse(layout, endStr); err == nil {
// 			endDatePtr = &endDate
// 		}
// 	}

// 	chartData, err := c.repo.GetDonutChartData(ctx, sprint, startDatePtr, endDatePtr)
// 	if err != nil {
// 		c.logger.Error("Failed to get chart data", zap.Error(err))
// 		http.Error(w, "Failed to process chart data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Normalize percentages
// 	for i := range chartData {
// 		totalProd := chartData[i].ProdSuccess + chartData[i].ProdError
// 		totalDev := chartData[i].DevSuccess + chartData[i].DevError

// 		if totalProd > 0 {
// 			chartData[i].ProdSuccess = (chartData[i].ProdSuccess / totalProd) * 100
// 			chartData[i].ProdError = (chartData[i].ProdError / totalProd) * 100
// 		}
// 		if totalDev > 0 {
// 			chartData[i].DevSuccess = (chartData[i].DevSuccess / totalDev) * 100
// 			chartData[i].DevError = (chartData[i].DevError / totalDev) * 100
// 		}
// 	}

//		w.Header().Set("Content-Type", "application/json")
//		if err := json.NewEncoder(w).Encode(chartData); err != nil {
//			c.logger.Error("Failed to encode response", zap.Error(err))
//			http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
//		}
//	}
func (c *ChartController) GetDonutChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	sprint := query.Get("sprint")
	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	fmt.Println("🎯 /analytics/donut called")
	fmt.Println("sprint =", sprint)
	fmt.Println("startDate =", startDateStr)
	fmt.Println("endDate =", endDateStr)

	var startDatePtr, endDatePtr *time.Time
	layout := "2006-01-02"

	if startDateStr != "" {
		if startDate, err := time.Parse(layout, startDateStr); err == nil {
			startDatePtr = &startDate
		}
	}
	if endDateStr != "" {
		if endDate, err := time.Parse(layout, endDateStr); err == nil {
			endDatePtr = &endDate
		}
	}

	chartData, err := c.repo.GetDonutChartData(ctx, sprint, startDatePtr, endDatePtr)
	if err != nil {
		c.logger.Error("Failed to get chart data", zap.Error(err))
		http.Error(w, "Failed to process chart data", http.StatusInternalServerError)
		return
	}

	// Normalize percentages
	for i := range chartData {
		totalProd := chartData[i].ProdSuccess + chartData[i].ProdError
		totalDev := chartData[i].DevSuccess + chartData[i].DevError

		if totalProd > 0 {
			chartData[i].ProdSuccess = (chartData[i].ProdSuccess / totalProd) * 100
			chartData[i].ProdError = (chartData[i].ProdError / totalProd) * 100
		}
		if totalDev > 0 {
			chartData[i].DevSuccess = (chartData[i].DevSuccess / totalDev) * 100
			chartData[i].DevError = (chartData[i].DevError / totalDev) * 100
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chartData); err != nil {
		c.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
	}
}
