package service

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password string, hash string) error
}
