package model

import (
	"time"
	"github.com/google/uuid"
)


type Question struct {
	QuestionID   uuid.UUID `json:"question_id" gorm:"type:uuid;primaryKey"`
	ArchiveDate  time.Time `json:"archive_date" gorm:"type:date;not null"`
	QuestionText string    `json:"question_text" gorm:"type:varchar(255);not null"`
	YesVotes     int       `json:"yes_votes"    gorm:"not null"`
	NoVotes      int       `json:"no_votes"     gorm:"not null"`
	TotalVotes   int       `json:"total_votes"  gorm:"not null"`
	CreatedBy    int       `json:"created_by"` // FK to users.user_id
	CreatedAt    time.Time `json:"created_at"   gorm:"autoCreateTime"`
}