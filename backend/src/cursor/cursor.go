package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Base struct {
	HasNext bool `json:"has_next"`
}

func (c *Base) NoMore() bool {
	if c == nil {
		return false
	}
	return !c.HasNext
}

func (c *Base) SetHasNext(v bool) {
	c.HasNext = v
}

// 各嵌入具体实现
// 否则go编译器 生成的转发函数 会进行dereference, 导致panic
//func (c *Base) IsNil() bool {
//	return c == nil
//}

type CursorBase interface {
	ToMap() map[string]any
	NoMore() bool
	SetHasNext(bool)
	IsNil() bool
}

type Cursor struct {
	Version int             `json:"v"`
	Payload json.RawMessage `json:"p"`
}

func Encode(payload any) (string, error) {
	if payload == nil {
		return "", nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cursor payload marshal failed: %w", err)
	}

	c := Cursor{
		Version: 1,
		Payload: data,
	}

	raw, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("cursor wrapper marshal failed: %w", err)
	}

	return base64.URLEncoding.EncodeToString(raw), nil
}

func Decode[T any](encoded string) (*T, error) {
	if encoded == "" {
		return nil, nil
	}
	raw, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("cursor decode failed: %w", err)
	}

	var c Cursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("cursor unmarshal failed: %w", err)
	}

	if c.Version != 1 {
		return nil, fmt.Errorf("unsupported cursor version: %d", c.Version)
	}
	var target T
	if err := json.Unmarshal(c.Payload, &target); err != nil {
		return nil, fmt.Errorf("cursor payload unmarshal failed: %w", err)
	}

	return &target, nil
}

func Update[T CursorBase](old T, count, pageSize int, generate func(lastIdx int, hasNext bool) T) T {
	if count == 0 {
		if old.IsNil() {
			var zero T
			return zero
		}
		old.SetHasNext(false)
		return old
	}

	hasNext := count > pageSize
	lastIdx := pageSize - 1
	if !hasNext && count > 0 {
		lastIdx = count - 1
	}

	return generate(lastIdx, hasNext)
}

func NoMore(c CursorBase) bool {
	if c == nil {
		return false
	}
	if c.IsNil() {
		return false
	}

	return c.NoMore()
}

func HashNext(c CursorBase) bool {
	if c == nil {
		return false
	}
	if c.IsNil() {
		return false
	}

	return !c.NoMore()
}

func ToMap(c CursorBase) map[string]any {
	if c == nil || c.IsNil() {
		return nil
	}

	return c.ToMap()
}
