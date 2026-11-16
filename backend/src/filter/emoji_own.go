package filter

var EmojiOwnFilterDeinition = FilterDefinition{
	Fields: []FilterField{
		{Column: "source_type", Op: Eq},
	},
}

type EmojiOwnFilter struct {
	SourceType int8
}

func (f *EmojiOwnFilter) IsNil() bool {
	return f == nil
}

func (f *EmojiOwnFilter) ToMap() map[string]any {
	v := map[string]any{}
	if f.SourceType != 0 {
		v["source_type"] = f.SourceType
	}

	return v
}
