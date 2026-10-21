package request

type ConvMembListBrief4Conv struct {
	ConvID string `query:"conv_id" validate:"required"`
}

type ConvMembListBriefBatch struct {
	ConvID    string   `json:"conv_id" validate:"required"`
	MemberIDs []string `json:"member_ids" validate:"required"`
}
