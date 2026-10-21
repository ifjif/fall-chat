package command

type RestPassword struct {
	ID       int64
	Username string
	Email    string
	OldPass  string
	NewPass  string
	Code     string
	CodeID   string
}

type UpdateProfile struct {
	ID         int64
	Username   string
	Email      string
	Nickname   string
	Password   string
	SearchFlag int
	Avatar     string
	Phone      string
}

type ForgetPassword struct {
	Username string
	Email    string
	NewPass  string
	Code     string
	CodeID   string
}
