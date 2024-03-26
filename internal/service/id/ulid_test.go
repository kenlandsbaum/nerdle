package id

import (
	"testing"
	"time"
)

func TestGetUlid(t *testing.T) {

	u1 := GetUlid()
	time.Sleep(time.Millisecond * 1)
	u2 := GetUlid()
	time.Sleep(time.Millisecond * 1)
	u3 := GetUlid()

	compare12 := u1.Compare(u2)
	compare23 := u2.Compare(u3)
	compare13 := u1.Compare(u3)

	if compare12 != -1 {
		t.Errorf("expected -1 but got %d\n", compare12)
	}
	if compare23 != -1 {
		t.Errorf("expected -1 but got %d\n", compare23)
	}
	if compare13 != -1 {
		t.Errorf("expected -1 but got %d\n", compare13)
	}
}
