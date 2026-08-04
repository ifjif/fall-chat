package errs

import "fmt"

type BizError struct {
	Code    int // 业务码
	Status  int
	Message string
}

func (be *BizError) Error() string {
	return fmt.Sprintf("[%d] %s", be.Status, be.Message)
}

func NewBizError(code int, status int, message string) *BizError {
	return &BizError{Status: status, Code: code, Message: message}
}
