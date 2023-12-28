package domain

type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type AuthService interface {
	Login(*LoginRequest) (string, int, error)
}
