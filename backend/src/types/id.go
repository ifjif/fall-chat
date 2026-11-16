package types

import (
	"fmt"
	"strconv"
	"zzc/fall-chat/src/model"
)

type ID int64

func (id *ID) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		var err error
		s, err = strconv.Unquote(string(data))
		if err != nil {
			return fmt.Errorf("invalid id: %s", string(data))
		}
	}

	if s == "" {
		return nil
	}

	v, err := model.DecodeID(s)
	if err != nil {
		return fmt.Errorf("invalid id: %s", string(data))
	}

	*id = ID(v)
	return nil
}

func (id ID) Int64() int64 {
	return int64(id)
}
