package request

type ListIMMessage struct {
	ConvID  string `json:"conv_id"`
	Next    string `json:"next"`
	SortKey string `json:"sort_key"`
}
type PatchListIMMessage struct {
	ConvID string `json:"conv_id"`
	MaxSeq *int64 `json:"max_seq" validate:"required,min=0"`
	MinSeq *int64 `json:"min_seq" validate:"required,min=0"`
}
