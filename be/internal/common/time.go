package common

import "time"

func NowMillis() uint64 {
	return uint64(time.Now().UnixMilli())
}
