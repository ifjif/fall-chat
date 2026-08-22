package resp

import (
	"zzc/fall-chat/src/errs"

	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Responser struct {
	ctx     fiber.Ctx
	message string
	res     *Response
	err     *errs.BizError
}

func New(ctx fiber.Ctx) *Responser {
	return &Responser{ctx: ctx}
}

func (r *Responser) OK() *Responser {
	r.res = &Response{}
	r.res.Code = 200
	return r
}

func (r *Responser) Err(status int) *Responser {
	r.err = &errs.BizError{}
	r.err.Status = status
	r.err.Code = status
	return r
}

func (r *Responser) ErrCode(code int) *Responser {
	r.err.Code = code
	return r
}

func (r *Responser) Message(message string) *Responser {
	r.message = message
	return r
}

func (r *Responser) Failure() error {
	r.err.Message = r.message
	return r.err
}

func (r *Responser) JSON(data any) error {
	r.res.Data = data
	r.res.Message = r.message
	return r.ctx.Status(r.res.Code).JSON(r.res)
}
