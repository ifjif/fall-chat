package command

type RestPasswordCommand struct {
	ID       int64
	Username string
	Email    string
	OldPass  string
	NewPass  string
	Code     string
	CodeID   string
}

type UpdateProfileCommand struct {
	ID         int64
	Username   string
	Email      string
	Nickname   string
	Password   string
	SearchFlag int
	Avatar     string
	Phone      string
}

type ForgetPasswordCommand struct {
	Username string
	Email    string
	NewPass  string
	Code     string
	CodeID   string
}
