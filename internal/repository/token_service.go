package repository

type TokenService interface {
	Generate(userID string) (string, error)
	Validate(token string) (userID string, err error)
}
