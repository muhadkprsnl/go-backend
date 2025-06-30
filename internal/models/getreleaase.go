package models

// ReleaseResponse matches the format expected by the frontend
type ReleaseResponse struct {
	Version     string `json:"version"`
	Env         string `json:"env"`
	ReleaseDate string `json:"releaseDate"`
	CloseDate   string `json:"closeDate"`
	Status      string `json:"status"`
}
