package grpc

import (
	"context"

	taskpb "github.com/arkad0912/project-protos/proto/task"
	userpb "github.com/arkad0912/project-protos/proto/user"
	"github.com/arkad0912/tasks-service/internal/taskService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	svc        *taskService.TaskService
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *taskService.TaskService, uc userpb.UserServiceClient) *Handler {
	return &Handler{
		svc:        svc,
		userClient: uc,
	}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	// 1. Проверяем существование пользователя
	userResp, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		Id: req.GetUserId(),
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to verify user")
	}

	// 2. Создаём задачу
	task, err := h.svc.CreateTask(taskService.Task{
		Task:   req.GetTitle(),
		UserID: uint(req.GetUserId()),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(task.ID),
			Title:  task.Task,
			UserId: userResp.GetUser().GetId(), // Используем ID из ответа
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	task, err := h.svc.GetTaskByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(task.ID),
			Title:  task.Task,
			UserId: uint32(task.UserID),
		},
	}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	// Получаем существующую задачу
	existingTask, err := h.svc.GetTaskByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	// Обновляем только title (если он передан)
	if req.GetTitle() != "" {
		existingTask.Task = req.GetTitle()
	}

	// Сохраняем обновленную задачу
	updatedTask, err := h.svc.UpdateTaskByID(uint(req.GetId()), existingTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(updatedTask.ID),
			Title:  updatedTask.Task,
			UserId: uint32(updatedTask.UserID),
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	err := h.svc.DeleteTaskByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			UserId: uint32(t.UserID),
		})
	}

	return &taskpb.ListTasksResponse{
		Tasks:      pbTasks,
		TotalCount: uint32(len(tasks)),
	}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksByUserResponse, error) {
	tasks, err := h.svc.GetTasksByUserID(uint(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			UserId: uint32(t.UserID),
		})
	}

	return &taskpb.ListTasksByUserResponse{
		Tasks:      pbTasks,
		TotalCount: uint32(len(tasks)),
	}, nil
}
