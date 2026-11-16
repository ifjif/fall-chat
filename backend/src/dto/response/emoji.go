package response

import "time"

type EmojiPackResult struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CoverURL  string    `json:"cover_url"`
	OwnerID   string    `json:"owner_id"`
	Source    string    `json:"soruce"`     // 1:系统 2:用户
	PriceType string    `json:"price_type"` // 1免费 2付费
	Price     float64   `json:"price"`
	Status    string    `json:"status"` // 0:下架 1:草稿 2:上架 3:审核中
	CreatedAt time.Time `json:"created_at"`
}

type EmojiItemResult struct {
	ID        string `json:"id"`
	PackID    string `json:"pack_id"`
	MediaType string `json:"media_type"` // 1静态图 2动态图 3字体 4unicode
	MediaURL  string `json:"media_url"`
	FontName  string `json:"font_name"`
	ThumbURL  string `json:"thumb_url"`
	Width     uint16 `json:"width"`
	Height    uint16 `json:"height"`
	FileSize  uint   `json:"file_size"`
	SortNO    int    `json:"sort_no"`
}
