package entity

type CreateQuestionRequest struct {
	ArchiveDate  string `json:"archive_date"`  // Expected format: "2006-01-02"
	QuestionText string `json:"question_text"`
	YesVotes     int    `json:"yes_votes"`
	NoVotes      int    `json:"no_votes"`
	TotalVotes   int    `json:"total_votes"`
	CreatedBy    int    `json:"created_by"`
}