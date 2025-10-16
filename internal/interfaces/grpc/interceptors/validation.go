package interceptors

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/Roflan4eg/auth-serivce/internal/lib/validation"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	pb "github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
	"google.golang.org/grpc"
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
				formatValidationError(validationErr))
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
	err := validation.ValidateStruct(&validationReq)
	return err
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

func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errs []string
		for _, fieldErr := range validationErrors {
			errs = append(errs, formatFieldError(fieldErr))
		}
		return strings.Join(errs, "; ")
	}
	return err.Error()
}

func formatFieldError(fieldErr validator.FieldError) string {
	fieldName := strings.ToLower(fieldErr.Field())

	switch fieldErr.Tag() {
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fieldName, fieldErr.Param())
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fieldName, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fieldName, fieldErr.Param())
	case "strongPassword":
		return fmt.Sprintf("weak %s - must contain uppercase letters, numbers and special characters, at least 8 characters", fieldName)
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}
