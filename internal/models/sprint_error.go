package models

type SprintErrorRate struct {
	Name      string  `json:"name"`      // sprint name
	DevError  float64 `json:"devError"`  // %
	ProdError float64 `json:"prodError"` // %
}
