package util

func IF(f bool, a interface{}, b interface{}) interface{} {
	if f {
		return a
	}
	return b
}
