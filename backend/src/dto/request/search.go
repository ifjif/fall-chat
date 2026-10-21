package request

type UserOrGroup struct {
	Auth
	Name      string `json:"name" validate:"required"`
	NextCuror string `json:"next_cursor"`
}
