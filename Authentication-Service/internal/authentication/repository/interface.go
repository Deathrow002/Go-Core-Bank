package repository

import (
	"authentication-service/internal/authentication/models"
	"context"
)

type AuthenticationRepository interface {
	Create(ctx context.Context, authentication *models.Authentication) error
	GetByID(ctx context.Context, id string) (*models.Authentication, error)
	GetByUsername(ctx context.Context, username string) (*models.Authentication, error)
	GetByEmail(ctx context.Context, email string) (*models.Authentication, error)
	GetByCustomerID(ctx context.Context, customerID string) (*models.Authentication, error)
	Update(ctx context.Context, authentication *models.Authentication) error
	Delete(ctx context.Context, id string) error
}