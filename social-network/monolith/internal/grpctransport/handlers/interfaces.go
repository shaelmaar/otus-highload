package handlers

type AuthService interface {
	ValidateToken(string) (string, error)
}
