package service

import (
	"errors"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
)

var ErrInsufficientPrivileges = errors.New("insufficient privileges")

func AuthorizeByRole(identity *client.UserIdentity, requiredRole domain.UserRole) error {
	userRole := domain.UserRole(identity.UserRole)
	if userRole == requiredRole {
		return nil
	}

	if userRole == domain.Admin && requiredRole == domain.User {
		return nil
	}

	return ErrInsufficientPrivileges
}
