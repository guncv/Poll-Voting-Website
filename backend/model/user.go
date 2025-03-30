package model

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"password" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}