package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    "github.com/google/uuid"
    pb "github.com/manlikehenryy/grpc-task-manager/task" // replace with the actual path to the generated task.pb.go
)

type TaskServiceServer struct {
    pb.UnimplementedTaskServiceServer
    tasks map[string]*pb.Task
}

func NewTaskServiceServer() *TaskServiceServer {
    return &TaskServiceServer{
        tasks: make(map[string]*pb.Task),
    }
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
    id := uuid.New().String()
    task := &pb.Task{
        Id:          id,
        Title:       req.GetTitle(),
        Description: req.GetDescription(),
    }
    s.tasks[id] = task
    return &pb.TaskResponse{Task: task}, nil
}

func (s *TaskServiceServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
    task, exists := s.tasks[req.GetId()]
    if !exists {
        return nil, fmt.Errorf("task not found")
    }
    return &pb.TaskResponse{Task: task}, nil
}

func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
    task, exists := s.tasks[req.GetId()]
    if !exists {
        return nil, fmt.Errorf("task not found")
    }
    task.Title = req.GetTitle()
    task.Description = req.GetDescription()
    return &pb.TaskResponse{Task: task}, nil
}

func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
    _, exists := s.tasks[req.GetId()]
    if !exists {
        return nil, fmt.Errorf("task not found")
    }
    delete(s.tasks, req.GetId())
    return &pb.DeleteTaskResponse{Message: "Task deleted successfully"}, nil
}

func (s *TaskServiceServer) GetAllTasks(ctx context.Context, req *pb.GetAllTasksRequest) (*pb.GetAllTasksResponse, error) {
    tasks := make([]*pb.Task, 0, len(s.tasks))
    for _, task := range s.tasks {
        tasks = append(tasks, task)
    }
    return &pb.GetAllTasksResponse{Tasks: tasks}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterTaskServiceServer(grpcServer, NewTaskServiceServer())

    log.Println("gRPC server is running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
