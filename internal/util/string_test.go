package util

import "testing"

func TestInt64SliceToString(t *testing.T) {
	if r := Int64SliceToString([]int64{1, 2, 3, 4, 5}, ","); r != "1,2,3,4,5" {
		t.Error(r)
		t.FailNow()
	}
}
func TestInt64SliceToString2(t *testing.T) {
	if r := Int64SliceToString([]int64{}, ","); r != "" {
		t.Error(r)
		t.FailNow()
	}
}
func TestInt64SliceToString3(t *testing.T) {
	if r := Int64SliceToString(nil, ","); r != "" {
		t.Error(r)
		t.FailNow()
	}
}
func TestUint64SliceToString(t *testing.T) {
	if r := Int64SliceToString([]int64{1111111111111111111, 222222222222222222, 33333333333333333, 44444444444444, 5555555555555555}, "#"); r != "1111111111111111111#222222222222222222#33333333333333333#44444444444444#5555555555555555" {
		t.Error(r)
		t.FailNow()
	}
}
func TestUint64SliceToString2(t *testing.T) {
	if r := Int64SliceToString([]int64{}, ","); r != "" {
		t.Error(r)
		t.FailNow()
	}
}
func TestUint64SliceToString3(t *testing.T) {
	if r := Int64SliceToString(nil, ","); r != "" {
		t.Error(r)
		t.FailNow()
	}
}
func TestStringSliceToString(t *testing.T) {
	if r := StringSliceToString([]string{"111", "222", "333", "444"}, "-"); r != "111-222-333-444" {
		t.Error(r)
		t.FailNow()
	}
}
