package util

import (
	"bytes"
	"strconv"
)

func Int64SliceToString(nums []int64, sep string) string {
	if nums == nil || len(nums) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(nums[0], 10))
	for i := 1; i < len(nums); i++ {
		buf.WriteString(sep)
		buf.WriteString(strconv.FormatInt(nums[i], 10))
	}
	return buf.String()
}
func Uint64SliceToString(nums []uint64, sep string) string {
	if nums == nil || len(nums) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatUint(nums[0], 10))
	for i := 1; i < len(nums); i++ {
		buf.WriteString(sep)
		buf.WriteString(strconv.FormatUint(nums[i], 10))
	}
	return buf.String()
}
func StringSliceToString(strings []string, sep string) string {
	if strings == nil || len(strings) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(strings[0])
	for i := 1; i < len(strings); i++ {
		buf.WriteString(sep)
		buf.WriteString(strings[i])
	}
	return buf.String()
}
