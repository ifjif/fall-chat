package request

// 请求值 需要携带认证信息，从middleware中获取
type Auth struct {
	ID       int64  `json:"-"`
	Username string `json:"-"`
}
