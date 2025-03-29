package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/entity"
)

// CreateQuestion handles POST /question
func (s *Server) CreateQuestion(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: CreateQuestion] Called")
	
	var req entity.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: CreateQuestion] Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: CreateQuestion] Request parsed for question:", req.QuestionText)

	archiveDate, err := time.Parse("2006-01-02", req.ArchiveDate)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: CreateQuestion] Invalid archive_date format:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid archive_date format, use YYYY-MM-DD"})
	}

	// Pass context to the service call if your service supports it.
	question, err := s.questionService.CreateQuestion(c.Context(), archiveDate, req.QuestionText, req.YesVotes, req.NoVotes, req.TotalVotes, req.CreatedBy)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: CreateQuestion] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	s.logger.InfoWithID(c.Context(), "[Controller: CreateQuestion] Question created successfully")
	return c.Status(fiber.StatusCreated).JSON(question)
}

// GetAllQuestions handles GET /question
func (s *Server) GetAllQuestions(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: GetAllQuestions] Called")
	
	// Pass context to the service call if supported.
	questions, err := s.questionService.GetAllQuestions(c.Context())
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: GetAllQuestions] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: GetAllQuestions] Retrieved questions successfully")
	return c.Status(fiber.StatusOK).JSON(questions)
}

// GetQuestion handles GET /question/:id
func (s *Server) GetQuestion(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: GetQuestion] Called")
	
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: GetQuestion] Invalid question ID:", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid question ID"})
	}

	q, err := s.questionService.GetQuestionByID(c.Context(), id)
	if err != nil {
		if err.Error() == "question not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: GetQuestion] Question not found with id:", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: GetQuestion] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: GetQuestion] Retrieved question with id:", id)
	return c.Status(fiber.StatusOK).JSON(q)
}

// DeleteQuestion handles DELETE /question/:id
func (s *Server) DeleteQuestion(c *fiber.Ctx) error {
	s.logger.InfoWithID(c.Context(), "[Controller: DeleteQuestion] Called")
	
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.logger.ErrorWithID(c.Context(), "[Controller: DeleteQuestion] Invalid question ID:", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid question ID"})
	}

	if err := s.questionService.DeleteQuestion(c.Context(), id); err != nil {
		if err.Error() == "question not found" {
			s.logger.ErrorWithID(c.Context(), "[Controller: DeleteQuestion] Question not found with id:", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
		}
		s.logger.ErrorWithID(c.Context(), "[Controller: DeleteQuestion] Service error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	s.logger.InfoWithID(c.Context(), "[Controller: DeleteQuestion] Question deleted successfully with id:", id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Question deleted successfully"})
}
