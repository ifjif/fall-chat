package request

type ConvMembListBrief4Conv struct {
	ConvID string `query:"conv_id" validate:"required"`
}

type ConvMembListBriefBatch struct {
	ConvID    string   `json:"conv_id" validate:"required"`
	MemberIDs []string `json:"member_ids" validate:"required"`
}

type ConvMembCommitRead struct {
	MsgID  string `json:"msg_id" validate:"required"`
	MsgSeq int64  `json:"msg_seq" validate:"gt=0"`
	ConvID string `json:"conv_id" validate:"required"`
}
