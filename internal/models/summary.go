package models

type SummaryData struct {
	TotalBugs   int `json:"totalBugs"`
	SuccessRate int `json:"successRate"`
	ErrorRate   int `json:"errorRate"`
	Delays      int `json:"delays"`
}

type SummaryResponse struct {
	Prod SummaryData `json:"prod"`
	Dev  SummaryData `json:"dev"`
}
