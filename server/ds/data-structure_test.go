package ds

import (
	"log"
	"testing"
)

// type addTest struct {
// 	arg, expected interface{}
// }

type TestPiece struct {
	Val int `json:"val"`
	Tes int `json:"tes"`
}

func TestPush(t *testing.T) {
	newList := Stack{}
	newList.Push(TestPiece{Val: 5, Tes: 7})
	newList.Push(6)
	newList.Push(2)
	newList.Push(3)
	marshal, err := newList.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	arr := newList.SerializeArray()
	newList.Pop()
	newList.Pop()
	newList.Pop()
	arr2 := newList.SerializeArray()
	log.Println(arr, arr2)
	err = newList.UnmarshalJSON(marshal)
	if err != nil {
		log.Println(err)
	}
}
