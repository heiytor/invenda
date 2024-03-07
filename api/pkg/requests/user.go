package requests

type UserCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
