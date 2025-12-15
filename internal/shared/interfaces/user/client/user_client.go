package client

type UserClient interface {
	VerifyCredentials(username, password string) error
}
