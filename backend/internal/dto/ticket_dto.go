package dto 

type TicketDTO struct {
	Title string `json:"title" validate:"required,min=5,max=30"`
	Description string `json:"decription" validate:"max:150"`
}