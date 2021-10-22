package util

import "time"

func GetCST8Time(t time.Time) time.Time {
	return t.In(time.FixedZone("CST", 8*3600))
}
