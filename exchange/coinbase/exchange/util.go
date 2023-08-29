package exchange

import (
	"encoding/json"
)

func unmarshalFloatString(bts []byte, f *float64) error {
	if len(bts) <= 2 {
		*f = 0.0
		return nil
	}
	return json.Unmarshal(bts[1:len(bts)-1], f)
}
