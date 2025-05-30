package main

import (
	"github.com/arkad0912/tasks-service/internal/database"

	"github.com/arkad0912/tasks-service/internal/taskService"

	"github.com/arkad0912/tasks-service/internal/transport/grpc"
)

func main() {
	database.InitDB()                                      // 1. Подключение к БД
	userRepo := taskService.NewUserRepository(database.DB) // 2. Репозиторий
	userService := taskService.NewUserService(userRepo)    // 3. Сервис
	userHandler := grpc.NewUserHandlers(userService)       // 4. gRPC обработчики
	grpc.RunServer(userHandler, ":50051")                  // 5. Запуск сервера
}
