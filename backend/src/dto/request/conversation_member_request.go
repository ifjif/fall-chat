package request

type ConvMembListBrief4ConvRequest struct {
	ConvID string `query:"conv_id" validate:"required"`
}

type ConvMembListBriefBatchRequest struct {
	ConvID    string   `json:"conv_id" validate:"required"`
	MemberIDs []string `json:"member_ids" validate:"required"`
}
