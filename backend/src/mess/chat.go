package mess

import (
	"encoding/json"
)

type Chat struct {
	Message json.RawMessage `json:"message"`
	ConvID  string          `json:"conv_id"`
	To      string          `json:"to"`
}
