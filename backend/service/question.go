package service

import (
	"context"
	"errors"
	"time"

	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/model"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"gorm.io/gorm"
)

// QuestionService defines business operations for questions.
type QuestionService interface {
	CreateQuestion(ctx context.Context, archiveDate time.Time, text string, yesVotes, noVotes, totalVotes, createdBy int) (model.Question, error)
	GetQuestionByID(ctx context.Context, id int) (model.Question, error)
	GetAllQuestions(ctx context.Context) ([]model.Question, error)
	DeleteQuestion(ctx context.Context, id int) error
}

type questionService struct {
	repo repository.QuestionRepository
	log  log.LoggerInterface
}

// NewQuestionService creates a new questionService with injected repository and logger.
func NewQuestionService(r repository.QuestionRepository, logger log.LoggerInterface) QuestionService {
	return &questionService{
		repo: r,
		log:  logger,
	}
}

func (qs *questionService) CreateQuestion(ctx context.Context, archiveDate time.Time, text string, yesVotes, noVotes, totalVotes, createdBy int) (model.Question, error) {
	qs.log.InfoWithID(ctx, "[Service: CreateQuestion] Called")
	q := model.Question{
		ArchiveDate:  archiveDate,
		QuestionText: text,
		YesVotes:     yesVotes,
		NoVotes:      noVotes,
		TotalVotes:   totalVotes,
		CreatedBy:    createdBy,
	}
	created, err := qs.repo.CreateQuestion(ctx, q)
	if err != nil {
		qs.log.ErrorWithID(ctx, "[Service: CreateQuestion] Error creating question:", err)
		return model.Question{}, err
	}
	qs.log.InfoWithID(ctx, "[Service: CreateQuestion] Question created with id:", created.QuestionID)
	return created, nil
}

func (qs *questionService) GetQuestionByID(ctx context.Context, id int) (model.Question, error) {
	qs.log.InfoWithID(ctx, "[Service: GetQuestionByID] Called for id:", id)
	q, err := qs.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			qs.log.ErrorWithID(ctx, "[Service: GetQuestionByID] Question not found with id:", id)
			return model.Question{}, errors.New("question not found")
		}
		qs.log.ErrorWithID(ctx, "[Service: GetQuestionByID] Error finding question:", err)
		return model.Question{}, err
	}
	qs.log.InfoWithID(ctx, "[Service: GetQuestionByID] Found question with id:", id)
	return q, nil
}

func (qs *questionService) GetAllQuestions(ctx context.Context) ([]model.Question, error) {
	qs.log.InfoWithID(ctx, "[Service: GetAllQuestions] Called")
	questions, err := qs.repo.FindAll(ctx)
	if err != nil {
		qs.log.ErrorWithID(ctx, "[Service: GetAllQuestions] Error retrieving questions:", err)
		return nil, err
	}
	qs.log.InfoWithID(ctx, "[Service: GetAllQuestions] Retrieved questions successfully")
	return questions, nil
}

func (qs *questionService) DeleteQuestion(ctx context.Context, id int) error {
	qs.log.InfoWithID(ctx, "[Service: DeleteQuestion] Called for id:", id)
	// Verify question exists.
	_, err := qs.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			qs.log.ErrorWithID(ctx, "[Service: DeleteQuestion] Question not found with id:", id)
			return errors.New("question not found")
		}
		qs.log.ErrorWithID(ctx, "[Service: DeleteQuestion] Error finding question:", err)
		return err
	}
	if err := qs.repo.DeleteQuestion(ctx, id); err != nil {
		qs.log.ErrorWithID(ctx, "[Service: DeleteQuestion] Error deleting question:", err)
		return err
	}
	qs.log.InfoWithID(ctx, "[Service: DeleteQuestion] Question deleted with id:", id)
	return nil
}
