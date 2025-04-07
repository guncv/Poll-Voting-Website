package model

type QuestionCache struct {
	QuestionID        string `json:"question_id"`
	UserID            string `json:"user_id"`
	Text              string `json:"text"`
	FirstChoice       string `json:"first_choice"`
	SecondChoice      string `json:"second_choice"`
	TotalParticipants int    `json:"total_participants"`
	FirstChoiceCount  int    `json:"first_choice_count"`
	SecondChoiceCount int    `json:"second_choice_count"`
	Milestones        string `json:"milestones"` // like "100:id1,150:id2"
	FollowUps         string `json:"follow_ups"` // optional
	GroupID           string `json:"group_id"`   // for grouping related questions
}