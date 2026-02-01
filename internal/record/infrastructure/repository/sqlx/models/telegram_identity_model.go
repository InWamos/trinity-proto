package models

import (
	"time"

	"github.com/google/uuid"
)

// TelegramIdentityModel represents the sqlx model for the telegram_user_identities table.
type TelegramIdentityModel struct {
	ID          uuid.UUID `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Username    string    `db:"username"`
	PhoneNumber string    `db:"phone_number"`
	Bio         string    `db:"bio"`
	AddedAt     time.Time `db:"added_at"`
	AddedByUser uuid.UUID `db:"added_by_user"`
}
