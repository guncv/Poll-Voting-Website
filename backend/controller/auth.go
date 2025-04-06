package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guncv/Poll-Voting-Website/backend/util"
)

// JWTMiddleware validates the access token and sets the user ID in the context.
func JWTMiddleware(c *fiber.Ctx) error {
	// Log that the middleware is called.
	fmt.Println("[JWTMiddleware] Called")

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		fmt.Println("[JWTMiddleware] Missing Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Expect the header to be in the format "Bearer <token>"
	var tokenStr string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
	if err != nil || tokenStr == "" {
		fmt.Println("[JWTMiddleware] Invalid token format")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	// Validate the access token.
	token, err := util.ValidateAccessToken(tokenStr)
	if err != nil || !token.Valid {
		fmt.Println("[JWTMiddleware] Invalid or expired token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("[JWTMiddleware] Invalid token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Assuming the user ID is stored in the "sub" claim.
	userID := claims["sub"].(string)
	fmt.Println("[JWTMiddleware] Setting userID in context:", userID)
	c.Locals("userID", userID)
	return c.Next()
}

// Login handles user login.
func (s *Server) Login(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Login] Called")

	// Parse credentials from request.
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Login] Request parsed for email:", req.Email)

	// Authenticate user.
	user, err := s.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Authentication failed:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Login] User authenticated:", req.Email)

	// Generate tokens.
	accessToken, err := util.GenerateAccessToken(user.UserID.String())
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Error generating access token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}
	refreshToken, err := util.GenerateRefreshToken(user.UserID.String())
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Login] Error generating refresh token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Login] Tokens generated for user:", req.Email)

	// Set the refresh token in an HttpOnly cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,              // Set true in production (HTTPS).
		SameSite: "Lax",             // Or "Strict" based on your security needs.
		Path:     "/refresh",        // Optionally restrict the path.
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	s.logger.InfoWithID(c.Context(), "[Controller: Login] Refresh token cookie set for user:", req.Email)

	// Return the access token in the JSON response.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
	})
}

// Refresh handles refreshing the access token using the refresh token stored in the cookie.
func (s *Server) Refresh(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Refresh] Called")

	// Get the refresh token from cookie.
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: Refresh] No refresh token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No refresh token provided"})
	}

	// Validate the refresh token.
	token, err := util.ValidateRefreshToken(refreshToken)
	if err != nil || !token.Valid {
		s.logger.ErrorWithID(c.Context(), "[Controller: Refresh] Invalid refresh token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.ErrorWithID(c.Context(), "[Controller: Refresh] Invalid refresh token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token claims"})
	}

	userID := claims["sub"].(string)
	s.logger.InfoWithID(c.Context(), "[Controller: Refresh] Refresh token validated for user:", userID)

	// Generate new access token.
	newAccessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: Refresh] Error generating new access token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate new access token"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Refresh] New access token generated for user:", userID)

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
		Secure:   true,
		SameSite: "Lax",
	})

	s.logger.InfoWithID(c.Context(), "[Controller: Logout] Refresh token cookie cleared")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}


// Profile returns the authenticated user's profile.
func (s *Server) Profile(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: Profile] Called")

	// Extract the userID from the context (set by the JWT middleware) as a string.
	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		s.logger.ErrorWithID(c.Context(), "[Controller: Profile] No userID found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: Profile] Retrieved userID from context:", userIDStr)

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
