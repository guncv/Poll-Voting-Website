package controller

import (
	"strconv"

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

// Login checks user credentials
func (s *Server) Login(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Login] Called")

	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: Login] Request parsed for email:", req.Email)
	u, err := s.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			s.logger.ErrorWithID(c.Context(), "[Controller: Login] Invalid credentials for email:", req.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: Login] Login successful for email:", req.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"user_id": u.UserID,
			"email":   u.Email,
		},
	})
}

// Profile (Placeholder) returns user #1
func (s *Server) Profile(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Profile] Called")
	u, err := s.userService.GetUserByID(c.Context(), 1)
	if err != nil {
		if err.Error() == "user not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: Profile] User not found for id: 1")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: Profile] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Profile] User found for id: 1")
	return c.Status(fiber.StatusOK).JSON(u)
}

// Logout - placeholder
func (s *Server) Logout(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Logout] Called")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful (placeholder)"})
}

// GetUser by ID
func (s *Server) GetUser(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: GetUser] Called")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: GetUser] Invalid user ID:", idParam)
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
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: DeleteUser] Invalid user ID:", idParam)
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
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: UpdateUser] Invalid user ID:", idParam)
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
