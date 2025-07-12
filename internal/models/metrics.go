package models

type SprintMetric struct {
	Sprint         string  `bson:"sprint" json:"sprint"`
	Environment    string  `bson:"environment" json:"environment"`
	TotalBugs      int     `bson:"totalBugs" json:"totalBugs"`
	TotalTestCases int     `bson:"totalTestCases" json:"totalTestCases"`
	Passed         int     `bson:"passed" json:"passed"`
	Failed         int     `bson:"failed" json:"failed"`
	SuccessRate    float64 `bson:"successRate" json:"successRate"`
}
