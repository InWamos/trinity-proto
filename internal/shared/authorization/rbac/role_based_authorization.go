package rbac

import (
	"errors"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/internal/user/domain"
)

var ErrInsufficientPrivileges = errors.New("insufficient privileges")

func AuthorizeByRole(identity *client.UserIdentity, requiredRole domain.Role) error {
	userRole := domain.Role(identity.UserRole)
	if userRole == requiredRole {
		return nil
	}

	if userRole == domain.RoleAdmin && requiredRole == domain.RoleUser {
		return nil
	}

	return ErrInsufficientPrivileges
}
