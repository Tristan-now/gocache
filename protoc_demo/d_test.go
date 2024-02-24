package protoc_demo

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func Test1(t *testing.T) {
	test := &Student{
		Name:   "geektutu",
		Male:   true,
		Scores: []int32{22, 1, 333},
	}
	fmt.Println(test)
	data, _ := proto.Marshal(test)
	fmt.Println(data)
	s1 := new(Student)
	proto.Unmarshal(data, s1)
	fmt.Println(s1)

}
