package model

import "time"

type User struct {
    UserID    int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"password" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
