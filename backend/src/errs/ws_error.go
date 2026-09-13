package errs

import "fmt"

type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (wse *WSError) Error() string {
	return fmt.Sprintf("WSError[%d] %s\n", wse.Code, wse.Message)
}

func NewWSError(code int, message string) *WSError {
	return &WSError{Code: code, Message: message}
}

var (
	ErrWSNotFound = NewWSError(404, "404 not found")
	ErrWSInternal = NewWSError(400, "websocket内部服务错误")
)
