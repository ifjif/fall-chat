package request

type ConvSingle struct {
	ConvID string `query:"conv_id" validate:"required"`
}
