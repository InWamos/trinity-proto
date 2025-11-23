package service

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct {
}

func (b *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	return string(bytes), err
}

func (b *BcryptPasswordHasher) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
