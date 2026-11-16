package filter

var EmojiPackFilterDefinition = FilterDefinition{
	Fields: []FilterField{
		{Column: "status", Op: Eq},
		{Column: "price_type", Op: Eq},
	},
}

type EmojiPackFilter struct {
	Status    *int8
	PriceType int8
}

func (f *EmojiPackFilter) IsNil() bool {
	return f == nil
}

func (f *EmojiPackFilter) ToMap() map[string]any {
	v := map[string]any{}
	if f.Status != nil {
		v["status"] = *f.Status
	}
	if f.PriceType != 0 {
		v["price_type"] = f.PriceType
	}

	return v
}
