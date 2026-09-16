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
	ErrWSNotFound                 = NewWSError(404, "404 not found")
	ErrWSInternal                 = NewWSError(400, "websocket内部服务错误")
	ErrWSClientIDSenderIDNotMatch = NewWSError(10000, "消息携带的发送id与连接的id不匹配")
	ErrWSChatConvType             = NewWSError(10001, "会话类型错误")
	ErrWSChatMessageType          = NewWSError(10002, "消息类型错误")
	ErrWSChatSenderNotExists      = NewWSError(10003, "消息发送者不存在")
	ErrWSChatReceiverNotExists    = NewWSError(10003, "消息接受者不存在")
)
