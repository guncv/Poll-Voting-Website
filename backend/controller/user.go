package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
    "github.com/guncv/Poll-Voting-Website/backend/entity"
)

// Register a new user
func (s *Server) Register(c *fiber.Ctx) error {
	var req entity.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := s.userService.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		s.logger.ErrorWithID(c.Context(), "Register error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": user.UserID,
		"email":   user.Email,
	})
}

// Login checks user credentials
func (s *Server) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := s.userService.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		s.logger.ErrorWithID(c.Context(), "Login error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Placeholder: In a real app, you'd create a session or set a JWT token.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"user_id": user.UserID,
			"email":   user.Email,
		},
	})
}

// Profile - a placeholder for returning the currently logged-in user's info
func (s *Server) Profile(c *fiber.Ctx) error {
	// In a real app, you'd parse a JWT from the headers, or check a session/cookie.
	// We'll pretend user #1 is logged in for now:
	user, err := s.userService.GetUserByID(1)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "Profile error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the user's info as an example
	return c.Status(fiber.StatusOK).JSON(user)
}

// Logout - a placeholder for clearing a session or token
func (s *Server) Logout(c *fiber.Ctx) error {
	// In a real app, you'd remove a session cookie or invalidate a JWT.
	// We'll just return a success message.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful (placeholder)",
	})
}

// GetUser fetches user by ID: GET /api/user/:id
func (s *Server) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := s.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "GetUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// DeleteUser removes a user by ID: DELETE /api/user/:id
func (s *Server) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = s.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "DeleteUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// UpdateUser modifies a user’s email/password: PUT /api/user/:id
func (s *Server) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req entity.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updatedUser, err := s.userService.UpdateUser(id, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "UpdateUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}
