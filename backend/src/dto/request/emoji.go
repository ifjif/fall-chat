package request

type CreateEmojiPack struct {
	Name      string  `json:"name" validate:"required"`
	CoverURL  string  `json:"cover_url"`
	Source    int8    `json:"source" validate:"oneof=1 2,required"`     // 1.系统 2.用户
	PriceType int8    `json:"price_type" validate:"oneof=1 2,required"` // 1免费 2付费
	Price     float64 `json:"price"`
}

type EmojiPackDetail struct {
	ID string `uri:"id" validate:""`
}

type ListEmojiPack struct {
	PriceType int8   `json:"price_type"`
	Status    *int8  `json:"status"`
	SortKey   string `json:"sort_key"`
	Next      string `json:"next"`
}

type CreateEmojiItem struct {
	MediaType int8   `json:"media_type"` // 1静态图 2动态图 3字体 4unicode
	MediaURL  string `json:"media_url"`
	FontName  string `json:"font_name"`
	ThumbURL  string `json:"thumb_url"`
	Width     uint16 `json:"width"`
	Height    uint16 `json:"height"`
	FileSize  uint   `json:"file_size"`
	SortNO    int    `json:"sort_no"`
}

type PublishEmoji struct {
	ID    string             `json:"id"`
	Items []*CreateEmojiItem `json:"items"`
}

type ListEmojiItem struct {
	PackID  string `json:"pack_id"`
	Next    string `json:"next"`
	SortKey string `json:"sort_key"` // sort_no_asc, sort_no_desc
}

type OwnListEmojiPack struct {
	SourceType int8   `json:"source_type"`
	Status     *int8  `json:"status"`
	Next       string `json:"next"`
}

type PublishEmojiVerifyCount struct {
	PackID    string `json:"pack_id" validate:"required"`
	ItemCount int    `json:"item_count" validate:"gt=0"`
}
