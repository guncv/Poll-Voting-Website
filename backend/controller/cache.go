package controller

import "github.com/gofiber/fiber/v2"

func (s *Server) getCache(c *fiber.Ctx) error {
	key := c.Params("key")
	val, err := s.cache.Get(key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if val == "" {
		return c.Status(fiber.StatusNotFound).SendString("Key not found")
	}
	return c.SendString(val)
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
