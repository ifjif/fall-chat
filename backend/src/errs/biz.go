package errs

import "fmt"

type BizError struct {
	Code    int // 业务码
	Status  int
	Message string
}

func (be *BizError) Error() string {
	return fmt.Sprintf("BizError[%d] %s", be.Status, be.Message)
}

func (be *BizError) ReMessage(message string, append bool) *BizError {
	m := be.Message
	if append {
		m = fmt.Sprintf("%s:%s", m, message)
	} else {
		m = message
	}
	nbe := &BizError{
		Code:    be.Code,
		Status:  be.Status,
		Message: m,
	}

	return nbe
}

// code = status
func NewStatusError(status int, message string) *BizError {
	return &BizError{Status: status, Code: status, Message: message}
}

func NewBizError(code, status int, message string) *BizError {
	return &BizError{Status: status, Code: code, Message: message}
}

func (b *BizError) Meger(err error) error {
	return fmt.Errorf("%w. internal: %w", b, err)
}
