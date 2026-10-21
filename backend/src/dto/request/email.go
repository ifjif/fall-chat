package request

type PassResetCode struct {
	Email string `json:"email" validate:"email"`
}
