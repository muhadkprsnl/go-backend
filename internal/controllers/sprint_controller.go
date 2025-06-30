package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SprintController struct {
	Collection *mongo.Collection
}

func NewSprintController(col *mongo.Collection) *SprintController {
	return &SprintController{Collection: col}
}

func (s *SprintController) GetSprints(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.TODO()
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$project", Value: bson.M{
			"sprint": bson.M{"$trim": bson.M{"input": "$sprint"}},
		}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$sprint"}}},
		bson.D{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := s.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	type SprintOption struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	var sprints []SprintOption
	for cursor.Next(ctx) {
		var result struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&result); err == nil && result.ID != "" {
			sprints = append(sprints, SprintOption{
				Value: result.ID,
				Label: result.ID,
			})
		}
	}

	if sprints == nil {
		sprints = []SprintOption{} // Ensure array, not null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sprints)
}
