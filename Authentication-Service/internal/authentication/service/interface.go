package service

import (
	"authentication-service/internal/authentication/models"
	"context"

	"github.com/google/uuid"
)

type AuthenticationService interface {
	CreateAuthentication(ctx context.Context, req CreateAuthenticationRequest) (*models.Authentication, error)
	ValidateLogin(ctx context.Context, username, password string) (*models.Authentication, error)
	GetByEmail(ctx context.Context, email string) (*models.Authentication, error)
	GetByCustomerID(ctx context.Context, customerID string) (*models.Authentication, error)
	UpdateAuthentication(ctx context.Context, id string, req UpdateAuthenticationRequest) (*models.Authentication, error)
	LockAccount(ctx context.Context, id string) error
	UnlockAccount(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error
	DeleteAuthentication(ctx context.Context, id string) error
}

type CreateAuthenticationRequest struct {
	CustomerID  uuid.UUID	`json:"customer_id" gorm:"type:uuid;not null;index"`
	Username    string 		`json:"username" validate:"required,min=3,max=50"`
	Password    string 		`json:"password" validate:"required,min=8"`
	Email       string 		`json:"email" validate:"required,email"`
	PhoneNumber string 		`json:"phone_number" validate:"required"`
	Role        string 		`json:"role" validate:"omitempty,oneof=admin support user"`
}

type UpdateAuthenticationRequest struct {
	Username    string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty" validate:"omitempty,oneof=admin support user"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}