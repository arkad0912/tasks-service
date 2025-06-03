package database

import (
	"log"

	// Импортируйте ваш пакет userService

	"github.com/arkad0912/tasks-service/internal/taskService"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=main port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Удаляем таблицу если существует
	if err := DB.Migrator().DropTable(&taskService.Task{}); err != nil {
		log.Println("Warning: could not drop table:", err)
	}

	// Создаем таблицу заново
	if err := DB.AutoMigrate(&taskService.Task{}); err != nil {
		log.Fatal("Failed to auto-migrate models:", err)
	}

	log.Println("Database migration completed successfully")
}
