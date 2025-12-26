package service

import (
	"testing"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/google/uuid"
)

func TestAuthorizeByRole(t *testing.T) {
	tests := []struct {
		name         string
		identity     *client.UserIdentity
		requiredRole domain.Role
		shouldPass   bool
	}{
		{
			name: "User accessing User-level endpoint",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: client.User,
			},
			requiredRole: domain.RoleUser,
			shouldPass:   true,
		},
		{
			name: "Admin accessing Admin-level endpoint",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: client.Admin,
			},
			requiredRole: domain.RoleAdmin,
			shouldPass:   true,
		},
		{
			name: "Admin accessing User-level endpoint",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: client.Admin,
			},
			requiredRole: domain.RoleUser,
			shouldPass:   true,
		},
		{
			name: "User accessing Admin-level endpoint",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: client.User,
			},
			requiredRole: domain.RoleAdmin,
			shouldPass:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AuthorizeByRole(tt.identity, tt.requiredRole)

			if tt.shouldPass {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if err != ErrInsufficientPrivileges {
					t.Errorf("expected ErrInsufficientPrivileges, got %v", err)
				}
			}
		})
	}
}

func TestAuthorizeByRoleEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		identity     *client.UserIdentity
		requiredRole domain.Role
		expectError  bool
	}{
		{
			name: "Admin role with empty required role",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: client.Admin,
			},
			requiredRole: "",
			expectError:  true,
		},
		{
			name: "Empty user role with RoleUser required role",
			identity: &client.UserIdentity{
				UserID:   uuid.New(),
				UserRole: "",
			},
			requiredRole: domain.RoleUser,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AuthorizeByRole(tt.identity, tt.requiredRole)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
