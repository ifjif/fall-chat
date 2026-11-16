package query

import "zzc/fall-chat/src/filter"

type EmojiPackDetail struct {
	ID        uint64
	Username  string
	PackageID int64
}

type ListEmojiPack struct {
	Filter  *filter.EmojiPackFilter
	Next    string
	SortKey string
}

type ListEmojiItem struct {
	PackID  int64
	Next    string
	SortKey string
}

type OwnListEmojiPack struct {
	UserID     uint64
	Next       string
	OwnFilter  *filter.EmojiOwnFilter
	PackFilter *filter.EmojiPackFilter
}

type PublishEmojiVerifyCount struct {
	UserID    uint64
	PackID    int64
	ItemCount int
}
