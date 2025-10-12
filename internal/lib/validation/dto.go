package validation

type CreateUserRequest struct {
	Email           string `validate:"required,email,min=5,max=255"`
	Password        string `validate:"required,strongPassword"`
	PasswordConfirm string `validate:"required,eqfield=Password"`
}

type GetUserRequest struct {
	ID string `validate:"required,uuid"`
}

type UpdateUserPasswordRequest struct {
	ID                 string `validate:"required,uuid"`
	OldPassword        string `validate:"required"`
	NewPassword        string `validate:"required,strongPassword"`
	NewPasswordConfirm string `validate:"required,eqfield=NewPassword"`
}

type UpdateUserEmailRequest struct {
	ID    string `validate:"required,uuid7"`
	Email string `validate:"required,email,min=5,max=255"`
}

type LoginRequest struct {
	Email    string `validate:"required,email,min=5,max=255"`
	Password string `validate:"required,strongPassword"`
}
