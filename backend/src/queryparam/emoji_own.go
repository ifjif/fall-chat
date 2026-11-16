package queryparam

import (
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/filter"
)

type EmojiOwnList struct {
	Owner      filter.ListFilterOptions
	Final      filter.ListFilterOptions
	Sort       cursor.ListSortOptions
	SourceType int8
}
