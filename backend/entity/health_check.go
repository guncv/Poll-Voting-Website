package entity

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}