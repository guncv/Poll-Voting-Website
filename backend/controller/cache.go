package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) getCache(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: getCache] Called")

	key := c.Params("key")
	date := time.Now().Format("2006-01-02")
	fullKey := "question:" + date + ":" + key

	result, err := s.cache.GetAllHash(fullKey)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: getCache] Redis error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Key not found"})
	}

	return c.JSON(result)
}

func (s *Server) setCache(c *fiber.Ctx) error {
	key := c.Params("key")
	value := c.FormValue("value")
	err := s.cache.Set(key, value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Saved")
}
