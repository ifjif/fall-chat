package request

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type LoginUser struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type ResetPassword struct {
	Auth
	Email   string `json:"email" validate:"email"`
	OldPass string `json:"old_pass" validate:"min=8,max=20"`
	NewPass string `json:"new_pass" validate:"min=8,max=20"`
	Code    string `json:"code" validate:"len=6"`
	CodeID  string `json:"code_id" validate:"required"`
}

type ForgetPassword struct {
	Username string `json:"username" validate:"min=3,max=50"`
	Email    string `json:"email" validate:"email"`
	NewPass  string `json:"new_pass" validate:"min=8,max=20"`
	Code     string `json:"code" validate:"len=6"`
	CodeID   string `json:"code_id" validate:"required"`
}

type UpdateProfile struct {
	Auth
	Email      string `json:"email" validate:"email"`
	Nickname   string `json:"nickname" validate:"min=3,max=50"`
	Password   string `json:"password" validate:"min=8,max=20"`
	SearchFlag int    `json:"search_flag" validate:"oneof=0 1"`
	Avatar     string `json:"avatar" validate:"required,max=255"`
	Phone      string `json:"phone" validate:"max=20"`
}
