package controller

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) HealthCheck(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[API: HealthCheck]: Called")
	response := s.healthCheckService.HealthCheck()
	s.logger.InfoWithID(c.Context(), "[API: HealthCheck]: Response: %v", response)
	return c.Status(fiber.StatusOK).JSON(response)
}