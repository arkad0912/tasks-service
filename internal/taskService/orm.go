package taskService

import (
	//"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task   string `gorm:"not null"`      // Наш сервер будет ожидать json c полем text
	IsDone bool   `gorm:"default:false"` // В GO используем CamelCase, в Json - snake
	UserID uint   `gorm:"not null"`      // Связь с пользователем
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
