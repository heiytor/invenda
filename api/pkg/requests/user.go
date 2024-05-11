package requests

type GetUser struct {
	ID string `param:"id" validate:"required|ulid"`
}

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

type CreateSession struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type UpdateSession struct {
	ID        string `param:"id" validate:"required"`
	Namespace string `json:"namespace"`
}
