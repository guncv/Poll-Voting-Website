package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/model"
	"gorm.io/gorm"
)

// QuestionRepository defines database operations for questions.
type QuestionRepository interface {
	CreateQuestion(ctx context.Context, q model.Question) (model.Question, error)
	FindByID(ctx context.Context, id int) (model.Question, error)
	FindAll(ctx context.Context) ([]model.Question, error)
	DeleteQuestion(ctx context.Context, id int) error
	FindLastArchivedQuestion(ctx context.Context) (model.Question, error)
}

type questionRepository struct {
	db  *gorm.DB
	log log.LoggerInterface
}

// NewQuestionRepository creates a new questionRepository with injected DB and logger.
func NewQuestionRepository(db *gorm.DB, logger log.LoggerInterface) QuestionRepository {
	return &questionRepository{
		db:  db,
		log: logger,
	}
}

func (qr *questionRepository) CreateQuestion(ctx context.Context, q model.Question) (model.Question, error) {
	qr.log.InfoWithID(ctx, "[Repository: CreateQuestion] Called for question:", q.QuestionText)
	if err := qr.db.Create(&q).Error; err != nil {
		qr.log.ErrorWithID(ctx, "[Repository: CreateQuestion] Error creating question:", err)
		return model.Question{}, err
	}
	qr.log.InfoWithID(ctx, "[Repository: CreateQuestion] Successfully created question with id:", q.QuestionID)
	return q, nil
}

func (qr *questionRepository) FindByID(ctx context.Context, id int) (model.Question, error) {
	qr.log.InfoWithID(ctx, "[Repository: FindByID] Called for question id:", id)
	var question model.Question
	if err := qr.db.First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			qr.log.ErrorWithID(ctx, "[Repository: FindByID] Question not found with id:", id)
			return model.Question{}, gorm.ErrRecordNotFound
		}
		qr.log.ErrorWithID(ctx, "[Repository: FindByID] Error finding question:", err)
		return model.Question{}, err
	}
	qr.log.InfoWithID(ctx, "[Repository: FindByID] Successfully found question with id:", id)
	return question, nil
}

func (qr *questionRepository) FindAll(ctx context.Context) ([]model.Question, error) {
	qr.log.InfoWithID(ctx, "[Repository: FindAll] Called")
	var questions []model.Question
	if err := qr.db.Find(&questions).Error; err != nil {
		qr.log.ErrorWithID(ctx, "[Repository: FindAll] Error retrieving questions:", err)
		return nil, err
	}
	qr.log.InfoWithID(ctx, "[Repository: FindAll] Successfully retrieved questions")
	return questions, nil
}

func (qr *questionRepository) DeleteQuestion(ctx context.Context, id int) error {
	qr.log.InfoWithID(ctx, "[Repository: DeleteQuestion] Called for question id:", id)
	if err := qr.db.Delete(&model.Question{}, id).Error; err != nil {
		qr.log.ErrorWithID(ctx, "[Repository: DeleteQuestion] Error deleting question:", err)
		return err
	}
	qr.log.InfoWithID(ctx, "[Repository: DeleteQuestion] Successfully deleted question with id:", id)
	return nil
}

func (qr *questionRepository) FindLastArchivedQuestion(ctx context.Context) (model.Question, error) {
	qr.log.InfoWithID(ctx, "[Repository: FindLastArchivedQuestion] Called")
    var q model.Question
    if err := qr.db.WithContext(ctx).
        Order("archive_date DESC").
        First(&q).Error; err != nil {
		qr.log.ErrorWithID(ctx, "[Repository: FindLastArchivedQuestion] Error finding last archived question:", err)
        return model.Question{}, err
    }
	if q.QuestionID == uuid.Nil {
		qr.log.ErrorWithID(ctx, "[Repository: FindLastArchivedQuestion] No archived question found")
		return model.Question{}, errors.New("no archived question found")
	}
	qr.log.InfoWithID(ctx, "[Repository: FindLastArchivedQuestion] Successfully found last archived question with id:", q.QuestionID)
    return q, nil
}