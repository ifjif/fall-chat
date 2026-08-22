package request

type PassResetCodeRequest struct {
	Email string `json:"email" validate:"email"`
}
