package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guncv/Poll-Voting-Website/backend/db"
	"github.com/guncv/Poll-Voting-Website/backend/entity"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/model"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"github.com/guncv/Poll-Voting-Website/backend/util"
	"gorm.io/gorm"
)

// QuestionService defines business operations for questions.
type IQuestionService interface {
	//DB question logic
	CreateQuestion(ctx context.Context,archiveDate time.Time,questionText string,firstChoice string,secondChoice string,totalParticipants int,firstChoiceCount int,secondChoiceCount int,createdBy uuid.UUID,) (model.Question, error)
	GetQuestionByID(ctx context.Context, id int) (model.Question, error)
	GetAllQuestions(ctx context.Context) ([]model.Question, error)
	DeleteQuestion(ctx context.Context, id int) error
	GetLastArchivedQuestion(ctx context.Context) (model.Question, error)

	// Redis vote logic
	VoteForQuestion(ctx context.Context, vote entity.VoteRequest) (entity.VoteResponse, error)

	// New Redis cache logic
	CreateQuestionCache(ctx context.Context, q entity.CreateQuestionCacheRequest) (string, error)
	GetQuestionCache(ctx context.Context, questionID string) (model.QuestionCache, error)
	DeleteQuestionCache(ctx context.Context, questionID string) error
	GetAllTodayQuestions(ctx context.Context) ([]model.QuestionCache, error)
}

type QuestionService struct {
	repo  repository.QuestionRepository
	cache db.CacheService
	log   log.LoggerInterface
}

// NewQuestionService creates a new questionService with injected repository and logger.
func NewQuestionService(r repository.QuestionRepository, cache db.CacheService, logger log.LoggerInterface) IQuestionService {
	return &QuestionService{
		repo:  r,
		cache: cache,
		log:   logger,
	}
}

func (qs *QuestionService) CreateQuestion(ctx context.Context,archiveDate time.Time,questionText string,firstChoice string,secondChoice string,totalParticipants int,firstChoiceCount int,secondChoiceCount int,createdBy uuid.UUID,) (model.Question, error) {

	qs.log.InfoWithID(ctx, "[Service: CreateQuestion] Called")

	q := model.Question{
		ArchiveDate:       archiveDate,
		QuestionText:      questionText,
		FirstChoice:       firstChoice,
		SecondChoice:      secondChoice,
		TotalParticipants: totalParticipants,
		FirstChoiceCount:  firstChoiceCount,
		SecondChoiceCount: secondChoiceCount,
		CreatedBy:         createdBy,
	}

	created, err := qs.repo.CreateQuestion(ctx, q)
	if err != nil {
		qs.log.ErrorWithID(ctx, "[Service: CreateQuestion] Error creating question:", err)
		return model.Question{}, err
	}

	qs.log.InfoWithID(ctx, "[Service: CreateQuestion] Question created with id:", created.QuestionID)
	return created, nil
}

func (qs *QuestionService) GetQuestionByID(ctx context.Context, id int) (model.Question, error) {
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

func (qs *QuestionService) GetAllQuestions(ctx context.Context) ([]model.Question, error) {
	qs.log.InfoWithID(ctx, "[Service: GetAllQuestions] Called")
	questions, err := qs.repo.FindAll(ctx)
	if err != nil {
		qs.log.ErrorWithID(ctx, "[Service: GetAllQuestions] Error retrieving questions:", err)
		return nil, err
	}
	qs.log.InfoWithID(ctx, "[Service: GetAllQuestions] Retrieved questions successfully")
	return questions, nil
}

func (qs *QuestionService) DeleteQuestion(ctx context.Context, id int) error {
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

func (qs *QuestionService) VoteForQuestion(ctx context.Context, vote entity.VoteRequest) (entity.VoteResponse, error) {
	date := util.TodayDate()
	qs.log.InfoWithID(ctx, "[Service: VoteForQuestion] Called for qid:", vote.QuestionID)

	voteKey := "voted:" + date + ":" + vote.QuestionID
	if voted, _ := qs.cache.IsSetMember(voteKey, vote.UserID); voted {
		qs.log.InfoWithID(ctx, "[Service: VoteForQuestion] User already voted")
		return entity.VoteResponse{AlreadyVoted: true, QuestionID: vote.QuestionID}, nil
	}

	qs.cache.AddSetMember(voteKey, vote.UserID)

	// Update vote counters
	field := "first_choice_count"
	if !vote.IsFirstChoice {
		field = "second_choice_count"
	}
	qs.cache.IncrementField("question:"+date+":"+vote.QuestionID, field)
	total := qs.cache.IncrementField("question:"+date+":"+vote.QuestionID, "total_participants")

	// Check milestone logic
	milestoneStr, _ := qs.cache.GetField("question:"+date+":"+vote.QuestionID, "milestones")
	revealedKey := "revealed:" + vote.QuestionID
	newlyRevealed := []string{}

	milestones := util.ParseMilestones(milestoneStr) // map[int]string

	for threshold, followUpID := range milestones {
		if total >= int64(threshold) {
			isRevealed, _ := qs.cache.IsSetMember(revealedKey, fmt.Sprint(threshold))
			if !isRevealed {
				qs.cache.AddSetMember(revealedKey, fmt.Sprint(threshold))
				newlyRevealed = append(newlyRevealed, followUpID)
			}
		}
	}

	first, _ := qs.cache.GetFieldInt("question:"+date+":"+vote.QuestionID, "first_choice_count")
	second, _ := qs.cache.GetFieldInt("question:"+date+":"+vote.QuestionID, "second_choice_count")

	return entity.VoteResponse{
		QuestionID:        vote.QuestionID,
		FirstChoiceCount:  first,
		SecondChoiceCount: second,
		TotalParticipants: int(total),
		NewlyRevealedIDs:  newlyRevealed,
		AlreadyVoted:      false,
	}, nil
}

func (qs *QuestionService) CreateQuestionCache(ctx context.Context, req entity.CreateQuestionCacheRequest) (string, error) {
	id := uuid.New().String()
	date := util.TodayDate()
	key := "question:" + date + ":" + id

	qs.log.InfoWithID(ctx, "[Service: CreateQuestionCache] Called for key:", key)

	// Build data as map for Redis
	data := map[string]string{
		"question_id":         id,
		"user_id":             req.UserID,
		"text":                req.Text,
		"first_choice":        req.FirstChoice,
		"second_choice":       req.SecondChoice,
		"first_choice_count":  "0",
		"second_choice_count": "0",
		"total_participants":  "0",
		"milestones":          req.Milestones,
		"follow_ups":          req.FollowUps,
		"group_id":            req.GroupID,
	}

	if err := qs.cache.SetHash(key, data); err != nil {
		qs.log.ErrorWithID(ctx, "[Service: CreateQuestionCache] Failed to store in Redis:", err)
		return "", err
	}

	if err := qs.cache.AddToSet("questions:"+date, id); err != nil {
		return "", err
	}

	return id, nil
}

func (qs *QuestionService) GetQuestionCache(ctx context.Context, questionID string) (model.QuestionCache, error) {
	date := util.TodayDate()
	key := "question:" + date + ":" + questionID
	qs.log.InfoWithID(ctx, "[Service: GetQuestionCache] Called for key:", key)

	data, err := qs.cache.GetAllHash(key)
	if err != nil {
		qs.log.ErrorWithID(ctx, "[Service: GetQuestionCache] Failed:", err)
		return model.QuestionCache{}, err
	}

	return model.QuestionCache{
		QuestionID:        data["question_id"],
		UserID:            data["user_id"],
		Text:              data["text"],
		FirstChoice:       data["first_choice"],
		SecondChoice:      data["second_choice"],
		FirstChoiceCount:  util.AtoiOrZero(data["first_choice_count"]),
		SecondChoiceCount: util.AtoiOrZero(data["second_choice_count"]),
		TotalParticipants: util.AtoiOrZero(data["total_participants"]),
		Milestones:        data["milestones"],
		FollowUps:         data["follow_ups"],
		GroupID:           data["group_id"],
	}, nil
}

func (qs *QuestionService) DeleteQuestionCache(ctx context.Context, questionID string) error {
	date := util.TodayDate()
	key := "question:" + date + ":" + questionID
	qs.log.InfoWithID(ctx, "[Service: DeleteQuestionCache] Deleting key:", key)
	return qs.cache.DeleteKey(key)
}

func (qs *QuestionService) GetAllTodayQuestions(ctx context.Context) ([]model.QuestionCache, error) {
	date := util.TodayDate()
	key := "questions:" + date
	qs.log.InfoWithID(ctx, "[Service: GetAllTodayQuestions] Listing from key:", key)

	ids, err := qs.cache.GetSetMembers(key)
	if err != nil {
		return nil, err
	}

	var result []model.QuestionCache
	for _, id := range ids {
		fullKey := "question:" + date + ":" + id
		data, err := qs.cache.GetAllHash(fullKey)
		if err != nil {
			qs.log.ErrorWithID(ctx, "[Service: GetAllTodayQuestions] Failed to fetch for key:", fullKey)
			continue
		}

		question := model.QuestionCache{
			QuestionID:        data["question_id"],
			UserID:            data["user_id"],
			Text:              data["text"],
			FirstChoice:       data["first_choice"],
			SecondChoice:      data["second_choice"],
			FirstChoiceCount:  util.AtoiOrZero(data["first_choice_count"]),
			SecondChoiceCount: util.AtoiOrZero(data["second_choice_count"]),
			TotalParticipants: util.AtoiOrZero(data["total_participants"]),
			Milestones:        data["milestones"],
			FollowUps:         data["follow_ups"],
			GroupID:           data["group_id"],
		}
		result = append(result, question)
	}

	return result, nil
}

func (qs *QuestionService) GetLastArchivedQuestion(ctx context.Context) (model.Question, error) {
    qs.log.InfoWithID(ctx, "[Service: GetLastArchivedQuestion] Called")
    q, err := qs.repo.FindLastArchivedQuestion(ctx)
    if err != nil {
        qs.log.ErrorWithID(ctx, "[Service: GetLastArchivedQuestion] Error retrieving last archived question:", err)
        return model.Question{}, err
    }
    qs.log.InfoWithID(ctx, "[Service: GetLastArchivedQuestion] Found last archived question with id:", q.QuestionID)
    return q, nil
}