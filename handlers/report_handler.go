package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type DonutChartData struct {
	Name        string  `json:"name"`
	ProdSuccess float64 `json:"prodSuccess"`
	ProdError   float64 `json:"prodError"`
	DevSuccess  float64 `json:"devSuccess"`
	DevError    float64 `json:"devError"`
}

type ReportHandler struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewReportHandler(collection *mongo.Collection) *ReportHandler {
	logger, _ := zap.NewProduction()
	return &ReportHandler{
		collection: collection,
		logger:     logger,
	}
}

func (h *ReportHandler) GetDonutChartData(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			h.logger.Error("Recovered from panic in GetDonutChartData", zap.Any("error", rec))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	enableCORS(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.logger.Info("Received request for donut chart data", zap.String("method", r.Method))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$project": bson.M{
				"developers": bson.M{
					"$concatArrays": bson.A{
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$developer1", ""}},
								"then": bson.A{
									bson.M{"name": "$developer1", "env": "dev", "passed": "$d1_passed", "failed": "$d1_failed"},
								},
								"else": bson.A{},
							},
						},
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$developer2", ""}},
								"then": bson.A{
									bson.M{"name": "$developer2", "env": "dev", "passed": "$d2_passed", "failed": "$d2_failed"},
								},
								"else": bson.A{},
							},
						},
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$P_developer1", ""}},
								"then": bson.A{
									bson.M{"name": "$P_developer1", "env": "prod", "passed": "$P_d1_passed", "failed": "$P_d1_failed"},
								},
								"else": bson.A{},
							},
						},
						bson.M{
							"$cond": bson.M{
								"if": bson.M{"$ne": bson.A{"$P_developer2", ""}},
								"then": bson.A{
									bson.M{"name": "$P_developer2", "env": "prod", "passed": "$P_d2_passed", "failed": "$P_d2_failed"},
								},
								"else": bson.A{},
							},
						},
					},
				},
			},
		}, {"$unwind": "$developers"},
		{
			"$match": bson.M{
				"developers.name": bson.M{"$ne": ""},
			},
		},

		// {"$unwind": "$developers"},
		{
			"$match": bson.M{
				"developers.name": bson.M{
					"$ne": "",
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"name": "$developers.name",
					"env":  "$developers.env",
				},

				// {"$unwind": "$developers"},
				// {
				// 	"$group": bson.M{
				// 		"_id": bson.M{
				// 			"name": "$developers.name",
				// 			"env":  "$developers.env",
				// 		},
				"passed": bson.M{"$sum": "$developers.passed"},
				"failed": bson.M{"$sum": "$developers.failed"},
			},
		},
		// {
		// 	"$group": bson.M{
		// 		"_id": "$_id.name",
		// 		"stats": bson.M{
		// 			"$push": bson.M{
		// 				"env":    "$_id.env",
		// 				"passed": "$passed",
		// 				"failed": "$failed",
		// 			},
		// 		},
		// 	},
		// },
		{
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
		{
			"$match": bson.M{
				"_id": bson.M{
					"$ne": "",
				},
			},
		},
		{
			"$project": bson.M{
				"name": "$_id",
				"prodSuccess": bson.M{
					"$reduce": bson.M{
						"input":        "$stats",
						"initialValue": 0,
						"in": bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$eq": bson.A{"$$this.env", "prod"}},
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
								"if":   bson.M{"$eq": bson.A{"$$this.env", "prod"}},
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
								"if":   bson.M{"$eq": bson.A{"$$this.env", "dev"}},
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
								"if":   bson.M{"$eq": bson.A{"$$this.env", "dev"}},
								"then": bson.M{"$add": bson.A{"$$value", "$$this.failed"}},
								"else": "$$value",
							},
						},
					},
				},
			},
		},
	}

	cursor, err := h.collection.Aggregate(ctx, pipeline)
	if err != nil {
		h.logger.Error("Aggregation failed", zap.Error(err))
		http.Error(w, "Failed to aggregate data", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []DonutChartData
	if err := cursor.All(ctx, &results); err != nil {
		h.logger.Error("Cursor decoding failed", zap.Error(err))
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}

	for i := range results {
		totalProd := results[i].ProdSuccess + results[i].ProdError
		totalDev := results[i].DevSuccess + results[i].DevError

		if totalProd > 0 {
			results[i].ProdSuccess = (results[i].ProdSuccess / totalProd) * 100
			results[i].ProdError = (results[i].ProdError / totalProd) * 100
		}
		if totalDev > 0 {
			results[i].DevSuccess = (results[i].DevSuccess / totalDev) * 100
			results[i].DevError = (results[i].DevError / totalDev) * 100
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		h.logger.Error("JSON encoding failed", zap.Error(err))
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully returned donut data", zap.Int("count", len(results)))
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
