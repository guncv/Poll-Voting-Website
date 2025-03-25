package service

import "github.com/guncv/Poll-Voting-Website/backend/entity"

// HealthCheckService defines the interface for health check operations.
type HealthCheckService interface {
	HealthCheck() entity.HealthCheckResponse
}

// healthCheckService is the concrete implementation of HealthCheckService.
type healthCheckService struct{}

// NewHealthCheckService creates a new instance of HealthCheckService.
func NewHealthCheckService() HealthCheckService {
	return &healthCheckService{}
}

// HealthCheck returns a health check response.
func (s *healthCheckService) HealthCheck() entity.HealthCheckResponse {
	return entity.HealthCheckResponse{
		Status:  "ok",
		Message: "API is healthy!",
	}
}
