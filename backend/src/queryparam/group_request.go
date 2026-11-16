package queryparam

import (
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/filter"
)

type GroupRequestList struct {
	Filter filter.ListFilterOptions
	Sort   cursor.ListSortOptions
}
