package singleflight

import (
	"fmt"
	"testing"
)

var s = "init"
var x = 1

func TestDo(t *testing.T) {
	var g Group
	//v, err := g.Do("key", func() (interface{}, error) {
	//	return "bar", nil
	//})
	//
	//if v != "bar" || err != nil {
	//	t.Errorf("Do v = %v, error = %v", v, err)
	//}

	for i := 0; i < 100; i++ {
		go g.Do(s, func() (interface{}, error) {
			x++
			fmt.Println(x)
			return nil, nil
		})
	}

	//for i := 0; i < 100; i++ {
	//	g.Do(s, func() (interface{}, error) {
	//		x++
	//		println(x)
	//		return nil, nil
	//	})
	//}

}
