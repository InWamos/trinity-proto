package client

import "github.com/google/uuid"

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
)

type VerifyCredentialsResponse struct {
	UserID   uuid.UUID
	UserRole UserRole
}
