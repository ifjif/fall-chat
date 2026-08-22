package request

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type ResetPasswordRequest struct {
	AuthRequest
	Email   string `json:"email" validate:"email"`
	OldPass string `json:"old_pass" validate:"min=8,max=20"`
	NewPass string `json:"new_pass" validate:"min=8,max=20"`
	Code    string `json:"code" validate:"len=6"`
	CodeID  string `json:"code_id" validate:"required"`
}

type ForgetPasswordRequest struct {
	Username string `json:"username" validate:"min=3,max=50"`
	Email    string `json:"email" validate:"email"`
	NewPass  string `json:"new_pass" validate:"min=8,max=20"`
	Code     string `json:"code" validate:"len=6"`
	CodeID   string `json:"code_id" validate:"required"`
}

type UpdateProfileRequest struct {
	AuthRequest
	Email      string `json:"email" validate:"email"`
	Nickname   string `json:"nickname" validate:"min=3,max=50"`
	Password   string `json:"password" validate:"min=8,max=20"`
	SearchFlag int    `json:"search_flag" validate:"oneof=0 1"`
	Avatar     string `json:"avatar" validate:"required,max=255"`
	Phone      string `json:"phone" validate:"max=20"`
}

type SearchRequest struct {
	AuthRequest
	Name      string `json:"name" validate:"required"`
	NextCuror string `json:"next_cursor"`
}
