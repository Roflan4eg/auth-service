package handlers

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
	"github.com/Roflan4eg/auth-serivce/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthGRPCHandler struct {
	authService *services.AuthService
	pb.UnimplementedAuthServiceServer
}

func (h *AuthGRPCHandler) RegisterHandler(server *grpc.Server) {
	pb.RegisterAuthServiceServer(server, h)
}

func NewAuthGRPCHandler(authService *services.AuthService) *AuthGRPCHandler {
	return &AuthGRPCHandler{authService: authService}
}

func (h *AuthGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.SessionResponse, error) {
	user, err := h.authService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	resp := &pb.SessionResponse{AccessToken: user.AccessToken, RefreshToken: user.RefreshToken}
	return resp, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.SessionResponse, error) {
	ses, err := h.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	resp := &pb.SessionResponse{AccessToken: ses.AccessToken, RefreshToken: ses.RefreshToken}
	return resp, nil
}

func (h *AuthGRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authService.Logout(ctx, req.GetAccessToken()); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.SessionResponse, error) {
	ses, err := h.authService.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	resp := &pb.SessionResponse{RefreshToken: ses.RefreshToken, AccessToken: ses.AccessToken}
	return resp, nil
}

func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	validRes := h.authService.ValidateToken(ctx, req.GetToken())

	resp := &pb.ValidateTokenResponse{
		Valid:     validRes.Valid,
		UserId:    validRes.UserId,
		SessionId: validRes.SessionId,
		Error:     validRes.Error,
	}
	return resp, nil
}

//func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
//
//}
