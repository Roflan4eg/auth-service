package handlers

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
	"github.com/Roflan4eg/auth-serivce/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	userService *services.UserService
	pb.UnimplementedUserServiceServer
}

func (h *UserGRPCHandler) RegisterHandler(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, h)
}

func NewUserGRPCHandler(userService *services.UserService) *UserGRPCHandler {
	return &UserGRPCHandler{userService: userService}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.CreateUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: user.ID.String(), Email: user.Email, IsActive: user.IsActive}, nil
}

func (h *UserGRPCHandler) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: user.ID.String(), Email: user.Email, IsActive: user.IsActive}, nil
}

func (h *UserGRPCHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: user.ID.String(), Email: user.Email, IsActive: user.IsActive}, nil
}

func (h *UserGRPCHandler) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	err := h.userService.UpdateUserPassword(ctx, req.GetId(), req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
