package client

type Input struct {
	Name string `json:"name" binding:"required"`
}
