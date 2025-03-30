package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/entity"
)

// Register a new user
func (s *Server) Register(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Register] Called")
	
	var req entity.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Register] Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: Register] Request parsed for email:", req.Email)

	// Pass the request context to the service call.
	user, err := s.userService.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			s.logger.ErrorWithID(c.Context(), "[Controller: Register] User already exists for email:", req.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: Register] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: Register] User registered successfully for email:", req.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": user.UserID,
		"email":   user.Email,
	})
}

// GetUser by ID
func (s *Server) GetUser(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: GetUser] Called")
	
	// Get the user ID as string directly.
	id := c.Params("id")
	if id == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: GetUser] Missing user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	
	u, err := s.userService.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: GetUser] User not found with id:", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: GetUser] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: GetUser] User retrieved with id:", id)
	return c.Status(fiber.StatusOK).JSON(u)
}

// DeleteUser
func (s *Server) DeleteUser(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: DeleteUser] Called")
	
	// Get the user ID as string.
	id := c.Params("id")
	if id == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: DeleteUser] Missing user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	
	if err := s.userService.DeleteUser(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: DeleteUser] User not found with id:", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: DeleteUser] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: DeleteUser] User deleted with id:", id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

// UpdateUser
func (s *Server) UpdateUser(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: UpdateUser] Called")
	
	// Get the user ID as string.
	id := c.Params("id")
	if id == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: UpdateUser] Missing user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req entity.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: UpdateUser] Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: UpdateUser] Request parsed for user id:", id)
	u, err := s.userService.UpdateUser(c.Context(), id, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: UpdateUser] User not found with id:", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: UpdateUser] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: UpdateUser] User updated successfully with id:", id)
	return c.Status(fiber.StatusOK).JSON(u)
}
