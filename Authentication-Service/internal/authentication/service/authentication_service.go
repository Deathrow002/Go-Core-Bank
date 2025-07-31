package service

import (
	"authentication-service/internal/authentication/models"
	"authentication-service/internal/authentication/repository"
	"authentication-service/internal/external"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authenticationService struct {
    repo           repository.AuthenticationRepository
    customerClient *external.CustomerClient // Add customer client
}

func NewAuthenticationService(repo repository.AuthenticationRepository) AuthenticationService {
	return &authenticationService{
		repo: repo,
	}
}

// CreateAuthentication creates a new authentication record with password hashing
func (s *authenticationService) CreateAuthentication(ctx context.Context, req CreateAuthenticationRequest) (*models.Authentication, error) {
	// Validate business rules
	if err := s.validateCreateAuthenticationRequest(req); err != nil {
		return nil, err
	}

	// **Call Customer Service to validate customer exists**
    customerExists, err := s.customerClient.ValidateCustomer(ctx, req.CustomerID)
    if err != nil {
        return nil, fmt.Errorf("failed to validate customer: %w", err)
    }

	if !customerExists {
        return nil, errors.New("customer does not exist")
    }

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	authentication := &models.Authentication{
		CustomerID:   req.CustomerID,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Role:         req.Role,
		IsLocked:     false, // Default to unlocked
	}

	// Save to database
	if err := s.repo.Create(ctx, authentication); err != nil {
		return nil, err
	}

	// Don't return password hash
	authentication.PasswordHash = ""
	return authentication, nil
}

// ValidateLogin validates user credentials
func (s *authenticationService) ValidateLogin(ctx context.Context, username, password string) (*models.Authentication, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	// Get authentication record
	auth, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if account is locked
	if auth.IsLocked {
		return nil, errors.New("account is locked")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Don't return password hash
	auth.PasswordHash = ""
	return auth, nil
}

// GetByEmail retrieves authentication by email
func (s *authenticationService) GetByEmail(ctx context.Context, email string) (*models.Authentication, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	auth, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	auth.PasswordHash = ""
	return auth, nil
}

// GetByCustomerID retrieves authentication by customer ID
func (s *authenticationService) GetByCustomerID(ctx context.Context, customerID string) (*models.Authentication, error) {
	if customerID == "" {
		return nil, errors.New("customer_id is required")
	}

	auth, err := s.repo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	auth.PasswordHash = ""
	return auth, nil
}

// UpdateAuthentication updates an existing authentication record
func (s *authenticationService) UpdateAuthentication(ctx context.Context, id string, req UpdateAuthenticationRequest) (*models.Authentication, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get existing record
	auth, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Username != "" {
		auth.Username = req.Username
	}
	if req.Email != "" {
		auth.Email = req.Email
	}
	if req.PhoneNumber != "" {
		auth.PhoneNumber = req.PhoneNumber
	}
	if req.Role != "" {
		if !s.isValidRole(req.Role) {
			return nil, errors.New("role must be one of: admin, support, user")
		}
		auth.Role = req.Role
	}

	// Update in database
	if err := s.repo.Update(ctx, auth); err != nil {
		return nil, err
	}

	// Don't return password hash
	auth.PasswordHash = ""
	return auth, nil
}

// LockAccount locks a user account
func (s *authenticationService) LockAccount(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	auth, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	auth.IsLocked = true
	return s.repo.Update(ctx, auth)
}

// UnlockAccount unlocks a user account
func (s *authenticationService) UnlockAccount(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	auth, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	auth.IsLocked = false
	return s.repo.Update(ctx, auth)
}

// Helper methods
func (s *authenticationService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *authenticationService) validateCreateAuthenticationRequest(req CreateAuthenticationRequest) error {
	if req.CustomerID == uuid.Nil {
		return errors.New("customer_id is required")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	if req.Role == "" {
		req.Role = string(models.RoleTypeUser) // Default role
	}
	if !s.isValidRole(req.Role) {
		return errors.New("role must be one of: admin, support, user")
	}
	return nil
}

func (s *authenticationService) isValidRole(role string) bool {
	validRoles := []models.RoleType{
		models.RoleTypeAdmin,
		models.RoleTypeSupport,
		models.RoleTypeUser,
	}
	for _, validRole := range validRoles {
		if string(validRole) == role {
			return true
		}
	}
	return false
}

// ChangePassword changes the password for a given username, verifying the old password first
func (s *authenticationService) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if oldPassword == "" {
		return errors.New("old password is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	auth, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	auth.PasswordHash = hashedPassword

	if err := s.repo.Update(ctx, auth); err != nil {
		return err
	}
	return nil
}

// DeleteAuthentication deletes an authentication record by ID.
func (s *authenticationService) DeleteAuthentication(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}