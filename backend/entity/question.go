package entity

type CreateQuestionCacheRequest struct {
	Text         string `json:"text"`
	FirstChoice  string `json:"first_choice"`
	SecondChoice string `json:"second_choice"`
	Milestones   string `json:"milestones"` // like "100:id1,150:id2"
	FollowUps    string `json:"follow_ups"` // optional
	GroupID      string `json:"group_id"`   // optional
	UserID       string `json:"user_id"`    // Injected in controller from JWT
}

type QuestionCache struct {
	QuestionID        string `json:"question_id"`
	GroupID           string `json:"group_id"`
	FirstChoice       string `json:"first_choice"`
	FirstChoiceCount  int    `json:"first_choice_count"`
	FollowUps         string `json:"follow_ups"`
	Milestones        string `json:"milestones"`
	TotalParticipants int    `json:"total_participants"`
	SecondChoice      string `json:"second_choice"`
	SecondChoiceCount int    `json:"second_choice_count"`
	Text              string `json:"text"`
	UserID            string `json:"user_id"`
}

type CreateQuestionRequest struct {
	ArchiveDate       string `json:"archive_date"`
	QuestionText      string `json:"question_text"`
	FirstChoice       string `json:"first_choice"`
	SecondChoice      string `json:"second_choice"`
	TotalParticipants int    `json:"total_participants"`
	FirstChoiceCount  int    `json:"first_choice_count"`
	SecondChoiceCount int    `json:"second_choice_count"`
	CreatedBy         string `json:"created_by"`
}

type VoteRequest struct {
	UserID        string `json:"user_id"`
	QuestionID    string `json:"question_id"`
	IsFirstChoice bool   `json:"is_first_choice"`
}

type VoteResponse struct {
	QuestionID        string   `json:"question_id"`
	TotalParticipants int      `json:"total_participants"`
	FirstChoiceCount  int      `json:"first_choice_count"`
	SecondChoiceCount int      `json:"second_choice_count"`
	NewlyRevealedIDs  []string `json:"newly_revealed_ids"`
	AlreadyVoted      bool     `json:"already_voted"`
}
