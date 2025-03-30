package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guncv/Poll-Voting-Website/backend/util"
)

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Expect the header to be in the format "Bearer <token>"
	var tokenStr string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
	if err != nil || tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Validate token (this is a simplified example)
	token, err := util.ValidateAccessToken(tokenStr)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Assuming the user ID is stored in the "sub" claim.
	c.Locals("userID", claims["sub"].(string))
	return c.Next()
}

// Login handles user login.
func (s *Server) Login(c *fiber.Ctx) error {
	// Parse credentials from request (assume you have a LoginRequest entity)
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Authenticate user (this is your service call)
	user, err := s.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.UserID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}
	refreshToken, err := util.GenerateRefreshToken(user.UserID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
	}

	// Set the refresh token in an HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,              // set to true in production with HTTPS
		SameSite: "Lax",             // or "Strict" if appropriate
		Path:     "/refresh",        // restrict the path if you like
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	// Return the access token in the JSON response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func (s *Server) Refresh(c *fiber.Ctx) error {
	// Get the refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No refresh token provided"})
	}

	// Validate the refresh token
	token, err := util.ValidateRefreshToken(refreshToken)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token claims"})
	}

	userID := claims["sub"].(string)

	// Generate new access token
	newAccessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate new access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

// Logout clears the refresh token cookie.
func (s *Server) Logout(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Logout] Called")
	
	// Clear the refresh token cookie by setting its value to empty
	// and its expiration date to a time in the past.
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,    // Use true in production over HTTPS.
		SameSite: "Lax",   // Or "Strict" based on your security needs.
	})
	
	// Optionally, if you're tracking refresh tokens server-side, revoke the token here.
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}

func (s *Server) Profile(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Profile] Called")

	// Extract the userID from the context (set by the JWT middleware) as a string.
	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: Profile] No userID found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Call the service to get the user profile, passing the userID as a string.
	user, err := s.userService.GetUserByID(c.Context(), userIDStr)
	if err != nil {
		if err.Error() == "user not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: Profile] User not found for id:", userIDStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: Profile] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: Profile] Retrieved profile for user:", userIDStr)
	return c.Status(fiber.StatusOK).JSON(user)
}