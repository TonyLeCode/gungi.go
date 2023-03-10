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

func TestStackPush(t *testing.T) {
	newList := Stack{}
	newList.Push(TestPiece{Val: 5, Tes: 7})
	newList.Push(6)
	newList.Push(2)
	newList.Push(3)
	newList.Print()
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

func TestLinkedListPush(t *testing.T) {
	newList := LinkedList{}
	newList.Push(6)
	newList.Push(2)
	newList.Push(3)
	newList.Push(6)
	newList.Print()
	log.Println()
	newList.Remove(newList.Tail)
	newList.Print()
	log.Println()
	newList.Remove(newList.Head.Next)
	newList.Print()
}
