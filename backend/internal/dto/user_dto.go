package dto

type RegisterDTO struct {
	Name string `json:"name" validate:"required,min=5,max=15"`
	Email string `json:"email" validate:"email,required,min=8,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type LoginDTO struct {
	Email string `json:"email" validate:"email,required,min=8,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
