package common

import "encoding/json"

func ToJsonIgnoreErr(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
