package interceptors

import (
	"context"
	validation "github.com/Roflan4eg/auth-serivce/internal/lib/validation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
)

func Validation() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var validationErr error

		switch r := req.(type) {
		case *pb.CreateUserRequest:
			validationErr = validateCreateUserReq(r)
		case *pb.RegisterRequest:
			validationErr = validateRegisterReq(r)
		case *pb.LoginRequest:
			validationErr = validateLoginReq(r)
		case *pb.UpdateUserPasswordRequest:
			validationErr = validateUpdateUserPassReq(r)
		case *pb.GetUserRequest:
			validationErr = validateGetUserReq(r)
		}

		if validationErr != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"validation failed: %v",
				validationErr,
			)
		}

		return handler(ctx, req)
	}
}

func validateCreateUserReq(req *pb.CreateUserRequest) error {
	validationReq := validation.CreateUserRequest{
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
	return validation.ValidateStruct(&validationReq)
}

func validateUpdateUserPassReq(req *pb.UpdateUserPasswordRequest) error {
	validationReq := validation.UpdateUserPasswordRequest{
		ID:                 req.GetId(),
		OldPassword:        req.GetOldPassword(),
		NewPassword:        req.GetNewPassword(),
		NewPasswordConfirm: req.GetNewPasswordConfirm(),
	}
	return validation.ValidateStruct(&validationReq)
}

func validateGetUserReq(req *pb.GetUserRequest) error {
	validationReq := validation.GetUserRequest{
		ID: req.GetUserId(),
	}
	return validation.ValidateStruct(&validationReq)
}

func validateRegisterReq(req *pb.RegisterRequest) error {
	validationReq := validation.CreateUserRequest{
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
	return validation.ValidateStruct(&validationReq)
}

func validateLoginReq(req *pb.LoginRequest) error {
	validationReq := validation.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	return validation.ValidateStruct(&validationReq)
}
