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
	// Проверяем существование пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d not found: %v", req.UserId, err)
	}

	// Создаем задачу
	task, err := h.svc.CreateTask(taskService.Task{
		Task:   req.Title,
		UserID: req.UserId,
		IsDone: false,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(task.ID),
			UserId: task.UserID,
			Title:  task.Task,
			IsDone: task.IsDone,
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	// Реализация аналогична CreateTask
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbTasks []*taskpb.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Task,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.GetTasksByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbTasks []*taskpb.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Task,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.Task, error) {
	// Реализация аналогична CreateTask
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := h.svc.DeleteTaskByID(req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}
