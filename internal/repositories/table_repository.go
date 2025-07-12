package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/muhadkprsnl/go-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type TableRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewTableRepository(db *mongo.Database, logger *zap.Logger) *TableRepository {
	return &TableRepository{
		collection: db.Collection("Report_2"),
		logger:     logger,
	}
}

// Fetch all reports for a given environment
func (r *TableRepository) GetReportsByEnvironment(env string) ([]models.FormData, error) {
	filter := bson.M{
		"environment": env,
		"developer1":  bson.M{"$ne": ""},
		"developer2":  bson.M{"$ne": ""},
		"dueDate":     bson.M{"$gte": primitive.NewDateTimeFromTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))},
	}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		r.logger.Error("Failed to fetch reports", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var reports []models.FormData
	if err := cursor.All(context.TODO(), &reports); err != nil {
		r.logger.Error("Failed to decode reports", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// Update a report by ID
// func (r *TableRepository) UpdateReport(id primitive.ObjectID, report models.FormData) error {
// 	filter := bson.M{"_id": id}
// 	update := bson.M{
// 		"$set": bson.M{
// 			"environment":    report.Environment,
// 			"sprint":         report.Sprint,
// 			"version":        report.Version,
// 			"dueDate":        report.DueDate,
// 			"closeDate":      report.CloseDate,
// 			"developer1":     report.Developer1,
// 			"d1Passed":       report.D1Passed,
// 			"d1Failed":       report.D1Failed,
// 			"developer2":     report.Developer2,
// 			"d2Passed":       report.D2Passed,
// 			"d2Failed":       report.D2Failed,
// 			"totalTestCases": report.Totaltestcase,
// 			"totalBugs":      report.Totalbugs,
// 			"feature":        report.Feature,
// 			"createdAt":      time.Now(),
// 		},
// 	}

// 	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		r.logger.Error("Failed to update report", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }

// Update an existing report
// UpdateReport updates an existing report
func (r *TableRepository) UpdateReport(id primitive.ObjectID, updateData models.FormData) error {
	// Create filter and update documents
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"sprint":         updateData.Sprint,
			"version":        updateData.Version,
			"dueDate":        updateData.DueDate,
			"closeDate":      updateData.CloseDate,
			"totalTestCases": updateData.Totaltestcase,
			"totalBugs":      updateData.Totalbugs,
			"developer1":     updateData.Developer1,
			"d1Passed":       updateData.D1Passed,
			"d1Failed":       updateData.D1Failed,
			"developer2":     updateData.Developer2,
			"d2Passed":       updateData.D2Passed,
			"d2Failed":       updateData.D2Failed,
		},
	}

	// Execute the update
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		r.logger.Error("Failed to update report",
			zap.Error(err),
			zap.String("id", id.Hex()))
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", id.Hex())
	}

	return nil
}

// Delete a report
func (r *TableRepository) DeleteReport(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		r.logger.Error("Failed to delete report",
			zap.Error(err),
			zap.String("id", id.Hex()))
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with id %s", id.Hex())
	}

	return nil
}

// Delete a report by ID
// func (r *TableRepository) DeleteReport(id primitive.ObjectID) error {
// 	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
// 	if err != nil {
// 		r.logger.Error("Failed to delete report", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }
