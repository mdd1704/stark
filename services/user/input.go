package user

type Input struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
	Password string `json:"password" binding:"required"`
}
