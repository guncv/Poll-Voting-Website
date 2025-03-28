package repository

import (
	"errors"

	"github.com/guncv/Poll-Voting-Website/backend/model"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(q model.Question) (model.Question, error)
	FindByID(id int) (model.Question, error)
	FindAll() ([]model.Question, error)
	DeleteQuestion(id int) error
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (qr *questionRepository) CreateQuestion(q model.Question) (model.Question, error) {
	if err := qr.db.Create(&q).Error; err != nil {
		return model.Question{}, err
	}
	return q, nil
}

func (qr *questionRepository) FindByID(id int) (model.Question, error) {
	var question model.Question
	if err := qr.db.First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Question{}, gorm.ErrRecordNotFound
		}
		return model.Question{}, err
	}
	return question, nil
}

func (qr *questionRepository) FindAll() ([]model.Question, error) {
	var questions []model.Question
	if err := qr.db.Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (qr *questionRepository) DeleteQuestion(id int) error {
	if err := qr.db.Delete(&model.Question{}, id).Error; err != nil {
		return err
	}
	return nil
}
