package command

type CreateEmojiPack struct {
	UserID    uint64
	Username  string
	Name      string
	CoverURL  string
	Source    int8 // 1.系统 2.用户
	PriceType int8 // 1免费 2付费
	Price     float64
}

type CreateEmojiItem struct {
	MediaType int8 // 1静态图 2动态图 3字体 4unicode
	MediaURL  string
	FontName  string
	ThumbURL  string
	Width     uint16
	Height    uint16
	FileSize  uint
	SortNO    int
}

type PublishEmoji struct {
	UserID   uint64
	Username string
	PackID   int64
	Items    []*CreateEmojiItem
}
