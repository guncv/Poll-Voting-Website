package model

import (
	"time"
	"github.com/google/uuid"
)

type Question struct {
	QuestionID         uuid.UUID `json:"question_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ArchiveDate        time.Time `json:"archive_date" gorm:"type:date;not null"`
	QuestionText       string    `json:"question_text" gorm:"type:varchar(255);not null"`
	FirstChoice        string    `json:"first_choice"  gorm:"type:varchar(255);not null"`
	SecondChoice       string    `json:"second_choice" gorm:"type:varchar(255);not null"`
	TotalParticipants  int       `json:"total_participants" gorm:"not null;default:0"`
	FirstChoiceCount   int       `json:"first_choice_count" gorm:"not null;default:0"`
	SecondChoiceCount  int       `json:"second_choice_count" gorm:"not null;default:0"`
	CreatedBy          uuid.UUID `json:"created_by" gorm:"type:uuid;"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
}
