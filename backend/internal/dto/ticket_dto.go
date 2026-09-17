package dto 

type Plevel string

const (
	High   Plevel = "high"
	Middle Plevel = "middle"
	Low    Plevel = "low"
)

type TicketDTO struct {
	Title string `json:"title" validate:"required,min=5,max=30"`
	Description string `json:"decription" validate:"max:150"`
	Priority Plevel `json:"priority" validate:"required,oneof=high middle low"`
}