package repositories

import (
	"context"
	"local/qa-report/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ReportRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewReportRepository(db *mongo.Database, logger *zap.Logger) *ReportRepository {
	return &ReportRepository{
		collection: db.Collection("Report_2"),
		logger:     logger,
	}
}

func (r *ReportRepository) CreateReport(report models.FormData) (string, error) {
	report.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(context.TODO(), report)
	if err != nil {
		r.logger.Error("Failed to create report", zap.Error(err))
		return "", err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return insertedID.Hex(), nil
}

func (r *ReportRepository) GetAllReports() ([]models.FormData, error) {
	var reports []models.FormData

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		r.logger.Error("Failed to get reports", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &reports); err != nil {
		r.logger.Error("Failed to decode reports", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

func (r *ReportRepository) UpdateReport(id string, report models.FormData) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": report}

	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		r.logger.Error("Failed to update report",
			zap.String("id", id),
			zap.Error(err))
	}
	return err
}

func (r *ReportRepository) DeleteReport(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		r.logger.Error("Failed to delete report",
			zap.String("id", id),
			zap.Error(err))
	}
	return err
}

// Optional: Add this if you need raw aggregation access
func (r *ReportRepository) Aggregate(pipeline []bson.M) (*mongo.Cursor, error) {
	return r.collection.Aggregate(context.TODO(), pipeline)
}

func (r *ReportRepository) GetSummaryData(startDate, endDate time.Time, sprint string) (models.SummaryResponse, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999000000, time.UTC)

	envs := []string{"development", "production"}
	summary := models.SummaryResponse{}

	for _, env := range envs {
		filter := bson.M{
			"environment": env,
			"dueDate": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}
		if sprint != "" {
			filter["sprint"] = sprint
		}

		cursor, err := r.collection.Find(context.TODO(), filter)
		if err != nil {
			r.logger.Error("Failed to fetch summary data", zap.String("env", env), zap.Error(err))
			continue
		}
		defer cursor.Close(context.TODO())

		var totalBugs, totalPassed, totalFailed, delays, count int

		for cursor.Next(context.TODO()) {
			var report models.FormData
			if err := cursor.Decode(&report); err != nil {
				r.logger.Warn("Failed to decode report", zap.Error(err))
				continue
			}

			passed := report.D1Passed + report.D2Passed
			failed := report.D1Failed + report.D2Failed
			totalPassed += passed
			totalFailed += failed
			totalBugs += report.Totalbugs

			if !report.Feature && report.CloseDate.After(report.DueDate.Add(24*time.Hour)) {
				delays++
			}
			count++

			r.logger.Info("Decoded Report",
				zap.String("Environment", report.Environment),
				zap.Int("TotalBugs", report.Totalbugs),
				zap.Int("D1Passed", report.D1Passed),
				zap.Int("D1Failed", report.D1Failed),
				zap.Int("D2Passed", report.D2Passed),
				zap.Int("D2Failed", report.D2Failed),
				zap.Bool("Feature", report.Feature),
				zap.Time("DueDate", report.DueDate),
				zap.Time("CloseDate", report.CloseDate),
			)
		}

		successRate := 0
		errorRate := 0
		if totalPassed+totalFailed > 0 {
			successRate = (totalPassed * 100) / (totalPassed + totalFailed)
			errorRate = 100 - successRate
		}

		envData := models.SummaryData{
			TotalBugs:   totalBugs,
			SuccessRate: successRate,
			ErrorRate:   errorRate,
			Delays:      delays,
		}

		if env == "development" {
			summary.Dev = envData
		} else {
			summary.Prod = envData
		}
	}

	r.logger.Info("Final Summary Computed", zap.Any("summary", summary))
	return summary, nil
}

// SprintError rate
func (r *ReportRepository) GetSprintErrorRates(startDate, endDate time.Time) ([]models.SprintErrorRate, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999000000, time.UTC)

	filter := bson.M{
		"dueDate": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	type Accumulator struct {
		TotalBugs int
		Failed    int
	}

	sprintMap := make(map[string]map[string]*Accumulator) // sprint -> env -> Acc

	for cursor.Next(context.TODO()) {
		var report models.FormData
		if err := cursor.Decode(&report); err != nil {
			continue
		}
		sprint := report.Sprint
		env := report.Environment
		if _, ok := sprintMap[sprint]; !ok {
			sprintMap[sprint] = map[string]*Accumulator{
				"development": {},
				"production":  {},
			}
		}
		acc := sprintMap[sprint][env]
		acc.TotalBugs += report.Totalbugs
		acc.Failed += report.D1Failed + report.D2Failed
	}

	var results []models.SprintErrorRate
	for sprint, envs := range sprintMap {
		result := models.SprintErrorRate{Name: sprint}

		if dev := envs["development"]; dev.TotalBugs > 0 {
			result.DevError = float64(dev.Failed) * 100 / float64(dev.TotalBugs)
		}
		if prod := envs["production"]; prod.TotalBugs > 0 {
			result.ProdError = float64(prod.Failed) * 100 / float64(prod.TotalBugs)
		}

		results = append(results, result)
	}

	return results, nil
}

// Releasereport
// ReleaseResponse represents the format expected by the frontend
type ReleaseResponse struct {
	Version     string `json:"version"`
	Env         string `json:"env"`
	ReleaseDate string `json:"releaseDate"`
	CloseDate   string `json:"closeDate"`
	Status      string `json:"status"` // "on time" or "delayed"
}

// GetReportsWithFilters filters by sprint and/or date and formats release data
// func (r *ReportRepository) GetRecentFilteredReports(ctx context.Context, sprint string, startDate, endDate *time.Time, status string) ([]models.ReleaseResponse, error) {
// 	filter := bson.M{}

// 	if sprint != "" && sprint != "All" {
// 		filter["sprint"] = sprint
// 	}
// 	if startDate != nil && endDate != nil {
// 		filter["dueDate"] = bson.M{
// 			"$gte": *startDate,
// 			"$lte": *endDate,
// 		}
// 	}

// 	// Sort by dueDate descending and limit to 5
// 	findOptions := options.Find()
// 	findOptions.SetSort(bson.D{{Key: "dueDate", Value: -1}})
// 	findOptions.SetLimit(5)

// 	cursor, err := r.collection.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var reports []models.FormData
// 	if err := cursor.All(ctx, &reports); err != nil {
// 		return nil, err
// 	}

// 	var response []models.ReleaseResponse
// 	for _, report := range reports {
// 		currentStatus := "on time"
// 		if !report.Feature && report.CloseDate.After(report.DueDate.Add(24*time.Hour)) {
// 			currentStatus = "delayed"
// 		}

// 		// Apply status filter
// 		if status != "" && status != "all" && currentStatus != status {
// 			continue
// 		}

// 		response = append(response, models.ReleaseResponse{
// 			Version:     report.Version,
// 			Env:         report.Environment,
// 			ReleaseDate: report.DueDate.Format("2006-01-02"),
// 			CloseDate:   report.CloseDate.Format("2006-01-02"),
// 			Status:      currentStatus,
// 		})
// 	}

// 	return response, nil
// }

// func (r *ReportRepository) GetRecentFilteredReports(ctx context.Context, status string) ([]models.ReleaseResponse, error) {
// 	// Extract skip and limit from context
// 	skipVal, ok := ctx.Value("skip").(int64)
// 	if !ok {
// 		skipVal = 0
// 	}

// 	limitVal, ok := ctx.Value("limit").(int64)
// 	if !ok || limitVal == 0 {
// 		limitVal = 5
// 	}

// 	filter := bson.M{}

// 	// Apply status filter
// 	if status == "delayed" {
// 		filter["feature"] = false
// 		filter["$expr"] = bson.M{
// 			"$gt": []interface{}{
// 				"$closeDate",
// 				bson.M{"$add": []interface{}{"$dueDate", 86400000}}, // +1 day in ms
// 			},
// 		}
// 	} else if status == "on-time" {
// 		filter["$or"] = []bson.M{
// 			{"feature": true},
// 			{
// 				"closeDate": bson.M{
// 					"$lte": bson.M{"$add": []interface{}{"$dueDate", 86400000}}, // <= dueDate + 1 day
// 				},
// 			},
// 		}
// 	}

// 	opts := options.Find().
// 		SetSort(bson.D{{"dueDate", -1}}).
// 		SetSkip(skipVal).
// 		SetLimit(limitVal)

// 	cursor, err := r.collection.Find(ctx, filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var reports []models.FormData
// 	if err := cursor.All(ctx, &reports); err != nil {
// 		return []models.ReleaseResponse{}, err
// 	}

// 	var response []models.ReleaseResponse
// 	for _, report := range reports {
// 		status := "on time"
// 		if !report.Feature && report.CloseDate.After(report.DueDate.Add(24*time.Hour)) {
// 			status = "delayed"
// 		}

// 		response = append(response, models.ReleaseResponse{
// 			Version:     report.Version,
// 			Env:         report.Environment,
// 			ReleaseDate: report.DueDate.Format("2006-01-02"),
// 			CloseDate:   report.CloseDate.Format("2006-01-02"),
// 			Status:      status,
// 		})
// 	}

// 	return []models.ReleaseResponse{}, err

func (r *ReportRepository) GetRecentFilteredReports(ctx context.Context, status string) ([]models.ReleaseResponse, error) {
	skipVal, ok := ctx.Value("skip").(int64)
	if !ok {
		skipVal = 0
	}
	limitVal, ok := ctx.Value("limit").(int64)
	if !ok || limitVal == 0 {
		limitVal = 5
	}

	filter := bson.M{}

	if status == "delayed" {
		filter["feature"] = false
		filter["$expr"] = bson.M{
			"$gt": []interface{}{
				"$closeDate",
				bson.M{"$add": []interface{}{"$dueDate", 86400000}}, // +1 day in ms
			},
		}
	} else if status == "on-time" {
		filter["$or"] = []bson.M{
			{"feature": true},
			{
				"$expr": bson.M{
					"$lte": []interface{}{
						"$closeDate",
						bson.M{"$add": []interface{}{"$dueDate", 86400000}},
					},
				},
			},
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "dueDate", Value: -1}}).
		SetSkip(skipVal).
		SetLimit(limitVal)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []models.FormData
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	var response []models.ReleaseResponse
	for _, report := range reports {
		status := "on time"
		if !report.Feature && report.CloseDate.After(report.DueDate.Add(24*time.Hour)) {
			status = "delayed"
		}

		response = append(response, models.ReleaseResponse{
			Version:     report.Version,
			Env:         report.Environment,
			ReleaseDate: report.DueDate.Format("2006-01-02"),
			CloseDate:   report.CloseDate.Format("2006-01-02"),
			Status:      status,
		})
	}

	return response, nil
}
