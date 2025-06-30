package models

type DonutChartData struct {
	Name        string  `json:"name"`
	ProdSuccess float64 `json:"prodSuccess"`
	ProdError   float64 `json:"prodError"`
	DevSuccess  float64 `json:"devSuccess"`
	DevError    float64 `json:"devError"`
}
