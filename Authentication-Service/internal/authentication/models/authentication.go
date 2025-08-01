package models

import (
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	RoleTypeAdmin    RoleType = "admin"
	RoleTypeSupport  RoleType = "support"
	RoleTypeUser     RoleType = "user"
)

type Authentication struct {
	ID           		uuid.UUID	`json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID   		uuid.UUID 	`json:"customer_id" gorm:"type:uuid;not null;index"`
	Username     		string 		`json:"username" gorm:"type:varchar(50);unique;not null;index"`
	PasswordHash 		string 		`json:"password_hash" gorm:"type:varchar(255);not null"`
	Email        		string 		`json:"email" gorm:"type:varchar(100);unique;not null;index"`
	PhoneNumber 		string 		`json:"phone_number" gorm:"type:varchar(20);not null;index"`
	Role 	 			string 		`json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	IsLocked 			bool   		`json:"is_locked" gorm:"type:boolean;default:false"`
	CreatedAt   		time.Time 	`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   		time.Time 	`json:"updated_at" gorm:"autoUpdateTime"`
}