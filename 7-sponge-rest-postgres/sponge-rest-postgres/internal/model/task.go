package model

import (
	"time"
)

type Task struct {
	ID          uint64     `gorm:"column:id;type:int4;primary_key" json:"id"`
	Title       string     `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Status      string     `gorm:"column:status;type:varchar(20)" json:"status"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamptz" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamptz" json:"updatedAt"`
}

// TableName table name
func (m *Task) TableName() string {
	return "task"
}

// TaskColumnNames Whitelist for custom query fields to prevent sql injection attacks
var TaskColumnNames = map[string]bool{
	"id":          true,
	"title":       true,
	"description": true,
	"status":      true,
	"created_at":  true,
	"updated_at":  true,
}
