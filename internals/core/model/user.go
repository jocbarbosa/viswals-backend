package model

import (
	"time"
)

// User defines the user structure from the file
type User struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	Email        string    `gorm:"uniqueIndex;size:100" json:"email_address"`
	CreatedAt    time.Time `json:"created_at"`
	MergedAt     time.Time `json:"merged_at"`
	ParentUserID int       `json:"parent_user_id"`
}
