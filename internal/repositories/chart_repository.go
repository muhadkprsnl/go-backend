package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/muhadkprsnl/go-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ChartRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewChartRepository(db *mongo.Database, logger *zap.Logger) *ChartRepository {
	return &ChartRepository{
		collection: db.Collection("Report_2"), // Same collection as reports
		logger:     logger,
	}
}

func (r *ChartRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// func (r *ChartRepository) GetDonutChartData(ctx context.Context, sprint string, startDate, endDate *time.Time) ([]models.DonutChartData, error) {
// 	matchStage := bson.M{}

// 	if sprint != "" && sprint != "All" {
// 		matchStage["sprint"] = strings.TrimSpace(sprint)
// 	}

// 	if startDate != nil && endDate != nil {
// 		matchStage["dueDate"] = bson.M{
// 			"$gte": *startDate,
// 			"$lte": *endDate,
// 		}
// 	}

// 	fmt.Println("🔍 Final matchStage:", matchStage)

// 	// Add the main aggregation stages
// 	pipeline := []bson.M{}

// 	if len(matchStage) > 0 {
// 		pipeline = append(pipeline, bson.M{"$match": matchStage})
// 	}

// 	pipeline = append(pipeline,
// 		bson.M{
// 			"$project": bson.M{
// 				"developers": bson.M{
// 					"$concatArrays": bson.A{
// 						bson.M{
// 							"$cond": bson.M{
// 								"if": bson.M{"$ne": bson.A{"$developer1", ""}},
// 								"then": bson.A{
// 									bson.M{
// 										"name":   "$developer1",
// 										"env":    "$environment", // development or production
// 										"passed": bson.M{"$ifNull": bson.A{"$d1Passed", 0}},
// 										"failed": bson.M{"$ifNull": bson.A{"$d1Failed", 0}},
// 									},
// 								},
// 								"else": bson.A{},
// 							},
// 						},
// 						bson.M{
// 							"$cond": bson.M{
// 								"if": bson.M{"$ne": bson.A{"$developer2", ""}},
// 								"then": bson.A{
// 									bson.M{
// 										"name":   "$developer2",
// 										"env":    "$environment",
// 										"passed": bson.M{"$ifNull": bson.A{"$d2Passed", 0}},
// 										"failed": bson.M{"$ifNull": bson.A{"$d2Failed", 0}},
// 									},
// 								},
// 								"else": bson.A{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		bson.M{"$unwind": "$developers"},
// 		bson.M{"$match": bson.M{"developers.name": bson.M{"$ne": ""}}},

// 		bson.M{
// 			"$group": bson.M{
// 				"_id": bson.M{
// 					"name": "$developers.name",
// 					"env":  "$developers.env",
// 				},
// 				"passed": bson.M{"$sum": "$developers.passed"},
// 				"failed": bson.M{"$sum": "$developers.failed"},
// 			},
// 		},
// 		bson.M{
// 			"$group": bson.M{
// 				"_id": "$_id.name",
// 				"stats": bson.M{
// 					"$push": bson.M{
// 						"env":    "$_id.env",
// 						"passed": "$passed",
// 						"failed": "$failed",
// 					},
// 				},
// 			},
// 		},
// 		bson.M{"$match": bson.M{"_id": bson.M{"$ne": ""}}},
// 		bson.M{
// 			"$project": bson.M{
// 				"name": "$_id",
// 				"prodSuccess": bson.M{
// 					"$reduce": bson.M{
// 						"input":        "$stats",
// 						"initialValue": 0,
// 						"in": bson.M{
// 							"$cond": bson.M{
// 								"if":   bson.M{"$eq": bson.A{"$$this.env", "production"}},
// 								"then": bson.M{"$add": bson.A{"$$value", "$$this.passed"}},
// 								"else": "$$value",
// 							},
// 						},
// 					},
// 				},
// 				"prodError": bson.M{
// 					"$reduce": bson.M{
// 						"input":        "$stats",
// 						"initialValue": 0,
// 						"in": bson.M{
// 							"$cond": bson.M{
// 								"if":   bson.M{"$eq": bson.A{"$$this.env", "production"}},
// 								"then": bson.M{"$add": bson.A{"$$value", "$$this.failed"}},
// 								"else": "$$value",
// 							},
// 						},
// 					},
// 				},
// 				"devSuccess": bson.M{
// 					"$reduce": bson.M{
// 						"input":        "$stats",
// 						"initialValue": 0,
// 						"in": bson.M{
// 							"$cond": bson.M{
// 								"if":   bson.M{"$eq": bson.A{"$$this.env", "development"}},
// 								"then": bson.M{"$add": bson.A{"$$value", "$$this.passed"}},
// 								"else": "$$value",
// 							},
// 						},
// 					},
// 				},
// 				"devError": bson.M{
// 					"$reduce": bson.M{
// 						"input":        "$stats",
// 						"initialValue": 0,
// 						"in": bson.M{
// 							"$cond": bson.M{
// 								"if":   bson.M{"$eq": bson.A{"$$this.env", "development"}},
// 								"then": bson.M{"$add": bson.A{"$$value", "$$this.failed"}},
// 								"else": "$$value",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	)

// 	cursor, err := r.collection.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		r.logger.Error("Aggregation failed", zap.Error(err))
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []models.DonutChartData
// 	if err := cursor.All(ctx, &results); err != nil {
// 		r.logger.Error("Failed to decode results", zap.Error(err))
// 		return nil, err
// 	}

//		return results, nil
//	}
func (r *ChartRepository) GetDonutChartData(ctx context.Context, sprint string, startDate, endDate *time.Time) ([]models.DonutChartData, error) {
	matchStage := bson.M{}

	if sprint != "" && sprint != "All" {
		matchStage["sprint"] = strings.TrimSpace(sprint)
	}

	if startDate != nil && endDate != nil {
		matchStage["dueDate"] = bson.M{
			"$gte": *startDate,
			"$lte": *endDate,
		}
	}

	fmt.Println("🔍 Final matchStage:", matchStage)

	pipeline := []bson.M{}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	pipeline = append(pipeline,
		bson.M{
			"$project": bson.M{
				"normalizedSprint": bson.M{"$trim": bson.M{"input": "$sprint"}},
				"developers": bson.M{
					"$concatArrays": bson.A{
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$developer1", ""}},
								"then": bson.A{
									bson.M{
										"name":   "$developer1",
										"env":    bson.M{"$toLower": "$environment"},
										"passed": bson.M{"$ifNull": bson.A{"$d1Passed", 0}},
										"failed": bson.M{"$ifNull": bson.A{"$d1Failed", 0}},
									},
								},
								"else": bson.A{},
							},
						},
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$developer2", ""}},
								"then": bson.A{
									bson.M{
										"name":   "$developer2",
										"env":    bson.M{"$toLower": "$environment"},
										"passed": bson.M{"$ifNull": bson.A{"$d2Passed", 0}},
										"failed": bson.M{"$ifNull": bson.A{"$d2Failed", 0}},
									},
								},
								"else": bson.A{},
							},
						},
					},
				},
			},
		},
		bson.M{"$unwind": "$developers"},
		bson.M{"$match": bson.M{"developers.name": bson.M{"$ne": ""}}},

		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"name": "$developers.name",
					"env":  "$developers.env",
				},
				"passed": bson.M{"$sum": "$developers.passed"},
				"failed": bson.M{"$sum": "$developers.failed"},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$_id.name",
				"stats": bson.M{
					"$push": bson.M{
						"env":    "$_id.env",
						"passed": "$passed",
						"failed": "$failed",
					},
				},
			},
		},
		bson.M{"$match": bson.M{"_id": bson.M{"$ne": ""}}},
		bson.M{
			"$project": bson.M{
				"name": "$_id",
				"prodSuccess": bson.M{
					"$reduce": bson.M{
						"input":        "$stats",
						"initialValue": 0,
						"in": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$$this.env", "production"}},
								"then": bson.M{"$add": bson.A{"$$value", "$$this.passed"}},
								"else": "$$value",
							},
						},
					},
				},
				"prodError": bson.M{
					"$reduce": bson.M{
						"input":        "$stats",
						"initialValue": 0,
						"in": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$$this.env", "production"}},
								"then": bson.M{"$add": bson.A{"$$value", "$$this.failed"}},
								"else": "$$value",
							},
						},
					},
				},
				"devSuccess": bson.M{
					"$reduce": bson.M{
						"input":        "$stats",
						"initialValue": 0,
						"in": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$$this.env", "development"}},
								"then": bson.M{"$add": bson.A{"$$value", "$$this.passed"}},
								"else": "$$value",
							},
						},
					},
				},
				"devError": bson.M{
					"$reduce": bson.M{
						"input":        "$stats",
						"initialValue": 0,
						"in": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$$this.env", "development"}},
								"then": bson.M{"$add": bson.A{"$$value", "$$this.failed"}},
								"else": "$$value",
							},
						},
					},
				},
			},
		},
	)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error("Aggregation failed", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.DonutChartData
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Error("Failed to decode results", zap.Error(err))
		return nil, err
	}

	for _, d := range results {
		fmt.Printf("👤 %s → Dev: %.2f/%.2f, Prod: %.2f/%.2f\n",
			d.Name, d.DevSuccess, d.DevError, d.ProdSuccess, d.ProdError)
	}

	return results, nil
}
