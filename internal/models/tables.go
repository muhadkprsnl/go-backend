package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableData struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Environment    string             `bson:"environment" json:"environment"`
	Sprint         string             `bson:"sprint" json:"sprint"`
	Version        string             `bson:"version" json:"version"`
	DueDate        time.Time          `bson:"dueDate" json:"dueDate"`
	CloseDate      time.Time          `bson:"closeDate" json:"closeDate"`
	TotalTestCases int                `bson:"totalTestCases" json:"totalTestCases"`
	TotalBugs      int                `bson:"totalBugs" json:"totalBugs"`
	Developers     []Developer        `bson:"developers" json:"developers"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
}

type Developer struct {
	Name   string `bson:"name" json:"name"`
	Passed int    `bson:"passed" json:"passed"`
	Failed int    `bson:"failed" json:"failed"`
}
