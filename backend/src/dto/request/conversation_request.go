package request

type ConvSingleRequest struct {
	ConvID string `query:"conv_id" validate:"required"`
}
