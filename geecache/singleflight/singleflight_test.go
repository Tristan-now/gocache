package singleflight

import (
	"testing"
)

var testCast = map[int]int{
	1: 1,
	2: 2,
	3: 3,
}

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}
