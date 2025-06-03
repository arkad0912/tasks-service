package grpc

import (
	"fmt"
	"log"
	"net"

	userpb "github.com/arkad0912/project-protos/proto/user"
	"github.com/arkad0912/tasks-service/internal/taskService"
	"google.golang.org/grpc"
)

func RunGRPC(svc *taskService.TaskService, uc userpb.UserServiceClient) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	srv := grpc.NewServer()
	handler := NewHandler(svc, uc)
	taskpb.RegisterTaskServiceServer(srv, handler)

	log.Printf("Tasks gRPC server started on :50052")
	return srv.Serve(lis)
}
