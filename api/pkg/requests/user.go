package requests

type CreateUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required|password"`
}

type UpdateUser struct {
	Name     string `json:"name" validate:""`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type AuthUser struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}
