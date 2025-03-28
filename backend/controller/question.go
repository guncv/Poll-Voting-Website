package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guncv/Poll-Voting-Website/backend/entity"
)

// CreateQuestion handles POST /question
func (s *Server) CreateQuestion(c *fiber.Ctx) error {
	var req entity.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	archiveDate, err := time.Parse("2006-01-02", req.ArchiveDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid archive_date format, use YYYY-MM-DD"})
	}

	question, err := s.questionService.CreateQuestion(
		archiveDate,
		req.QuestionText,
		req.YesVotes,
		req.NoVotes,
		req.TotalVotes,
		req.CreatedBy,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(question)
}

// GetAllQuestions handles GET /question
func (s *Server) GetAllQuestions(c *fiber.Ctx) error {
	questions, err := s.questionService.GetAllQuestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(questions)
}

// GetQuestion handles GET /question/:id
func (s *Server) GetQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid question ID"})
	}

	q, err := s.questionService.GetQuestionByID(id)
	if err != nil {
		if err.Error() == "question not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(q)
}

// DeleteQuestion handles DELETE /question/:id
func (s *Server) DeleteQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid question ID"})
	}

	if err := s.questionService.DeleteQuestion(id); err != nil {
		if err.Error() == "question not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Question deleted successfully"})
}
