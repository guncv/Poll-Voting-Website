package service

import (
	"errors"
	"time"

	"github.com/guncv/Poll-Voting-Website/backend/model"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"gorm.io/gorm"
)

type QuestionService interface {
	CreateQuestion(archiveDate time.Time, text string, yesVotes, noVotes, totalVotes, createdBy int) (model.Question, error)
	GetQuestionByID(id int) (model.Question, error)
	GetAllQuestions() ([]model.Question, error)
	DeleteQuestion(id int) error
}

type questionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(r repository.QuestionRepository) QuestionService {
	return &questionService{repo: r}
}

func (qs *questionService) CreateQuestion(archiveDate time.Time, text string, yesVotes, noVotes, totalVotes, createdBy int) (model.Question, error) {
	q := model.Question{
		ArchiveDate:  archiveDate,
		QuestionText: text,
		YesVotes:     yesVotes,
		NoVotes:      noVotes,
		TotalVotes:   totalVotes,
		CreatedBy:    createdBy,
	}
	created, err := qs.repo.CreateQuestion(q)
	if err != nil {
		return model.Question{}, err
	}
	return created, nil
}

func (qs *questionService) GetQuestionByID(id int) (model.Question, error) {
	q, err := qs.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Question{}, errors.New("question not found")
		}
		return model.Question{}, err
	}
	return q, nil
}

func (qs *questionService) GetAllQuestions() ([]model.Question, error) {
	return qs.repo.FindAll()
}

func (qs *questionService) DeleteQuestion(id int) error {
	_, err := qs.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}
	return qs.repo.DeleteQuestion(id)
}
