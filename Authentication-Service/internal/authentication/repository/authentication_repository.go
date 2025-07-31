package repository

import (
	"authentication-service/internal/authentication/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{
		db: db,
	}
}

func (r *authenticationRepository) Create(ctx context.Context, authentication *models.Authentication) error {
	if err := r.db.WithContext(ctx).Create(authentication).Error; err != nil {
		// Check for email conflicts (only if email is changing)
		var existing models.Authentication
		if existing.Email != authentication.Email {
			var count int64
			if err := r.db.WithContext(ctx).Model(&models.Authentication{}).
				Where("email = ?", authentication.Email).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("email address is already in use")
			}
		}

		// Check for username conflicts (only if username is changing)
		if existing.Username != authentication.Username {
			var count int64
			if err := r.db.WithContext(ctx).Model(&models.Authentication{}).
				Where("username = ?", authentication.Username).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("username is already taken")
			}
		}
	}
	return nil
}

func (r *authenticationRepository) ValidateByEmail(ctx context.Context, email string) (*models.Authentication, error) {
	var authentication models.Authentication
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&authentication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("authentication not found")
		}
		return nil, err
	}
	return &authentication, nil
}

func (r *authenticationRepository) ValidateByCustomerID(ctx context.Context, customerID string) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Authentication{}).Where("customer_id = ?", customerID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no authentication found for this customer")
	}
	return nil
}

func (r *authenticationRepository) Update(ctx context.Context, authentication *models.Authentication) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, check if record exists
		var existing models.Authentication
		if err := tx.Where("id = ?", authentication.ID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("authentication record not found")
			}
			return err
		}

		// Check for email conflicts (only if email is changing)
		if existing.Email != authentication.Email {
			var count int64
			if err := tx.Model(&models.Authentication{}).
				Where("email = ?", authentication.Email).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("email address is already in use")
			}
		}

		// Check for username conflicts (only if username is changing)
		if existing.Username != authentication.Username {
			var count int64
			if err := tx.Model(&models.Authentication{}).
				Where("username = ?", authentication.Username).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errors.New("username is already taken")
			}
		}

		// Perform the update
		return tx.Model(&models.Authentication{}).Where("id = ?", authentication.ID).Updates(authentication).Error
	})
}

func (r *authenticationRepository) GetByID(ctx context.Context, id string) (*models.Authentication, error) {
	var auth models.Authentication
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("authentication not found")
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authenticationRepository) GetByEmail(ctx context.Context, email string) (*models.Authentication, error) {
	var auth models.Authentication
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("authentication not found")
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authenticationRepository) GetByCustomerID(ctx context.Context, customerID string) (*models.Authentication, error) {
	var auth models.Authentication
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("authentication not found")
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authenticationRepository) GetByUsername(ctx context.Context, username string) (*models.Authentication, error) {
	var auth models.Authentication
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("authentication not found")
		}
		return nil, err
	}
	return &auth, nil
}

func (r *authenticationRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&models.Authentication{}, "id = ?", id).Error; err != nil {
		return errors.New("failed to delete authentication")
	}
	return nil
}

