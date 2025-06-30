package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormData struct {
	ObjectID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID            string             `json:"id,omitempty"`
	Environment   string             `json:"environment" bson:"environment"`
	Sprint        string             `json:"sprint" bson:"sprint"`
	Version       string             `json:"version" bson:"version"`
	DueDate       time.Time          `json:"dueDate" bson:"dueDate"`
	Totaltestcase int                `json:"totalTestCases" bson:"totalTestCases"`
	Totalbugs     int                `json:"totalBugs" bson:"totalBugs"`
	Developer1    string             `json:"developer1" bson:"developer1"`
	D1Passed      int                `json:"d1Passed" bson:"d1Passed"`
	D1Failed      int                `json:"d1Failed" bson:"d1Failed"`
	Developer2    string             `json:"developer2" bson:"developer2"`
	D2Passed      int                `json:"d2Passed" bson:"d2Passed"`
	D2Failed      int                `json:"d2Failed" bson:"d2Failed"`
	Feature       bool               `json:"feature" bson:"feature"` // ✅ Added this
	CloseDate     time.Time          `json:"closeDate" bson:"closeDate"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
}
