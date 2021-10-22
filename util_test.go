package biligo

import "testing"

func TestAV2BV(t *testing.T) {
	if AV2BV(170001) != "BV17x411w7KC" {
		t.FailNow()
	}
}
func TestBV2AV(t *testing.T) {
	if BV2AV("BV17x411w7KC") != 170001 {
		t.FailNow()
	}
}
