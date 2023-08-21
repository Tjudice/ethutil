package coinbase

import "encoding/json"

type FloatStringWrapper struct {
	Val *float64
}

func (f *FloatStringWrapper) UnmarshalJSON(bts []byte) error {
	if len(bts) <= 2 {
		*f.Val = 0.0
		return nil
	}
	return json.Unmarshal(bts[1:len(bts)-1], f.Val)
}
