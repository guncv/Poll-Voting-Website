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

type CreateQuestionRequest struct {
	ArchiveDate  string `json:"archive_date"` // Expected format: "2006-01-02"
	QuestionText string `json:"question_text"`
	YesVotes     int    `json:"yes_votes"`
	NoVotes      int    `json:"no_votes"`
	TotalVotes   int    `json:"total_votes"`
	CreatedBy    int    `json:"created_by"`
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
