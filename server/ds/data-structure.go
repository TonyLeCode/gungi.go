package ds

import (
	"encoding/json"
	"fmt"
)

type StackNode struct {
	Value interface{}
	Next  *StackNode
	Prev  *StackNode
}

// Stack using Doubly Linked List
type Stack struct {
	Bottom *StackNode
	Top    *StackNode
	Length int
}

func (stack *Stack) Push(value interface{}) {
	newStack := &StackNode{Value: value}
	if stack.Bottom == nil {
		stack.Bottom = newStack
		stack.Top = newStack
	} else {
		stack.Top.Next = newStack
		newStack.Prev = stack.Top
		stack.Top = newStack
	}
	stack.Length++
}

func (stack *Stack) Pop() {
	if stack.Top == nil {
		return
	}
	stack.Top = stack.Top.Prev
	if stack.Top != nil {
		stack.Top.Next = nil
	} else {
		stack.Bottom = nil
	}
	stack.Length--
}

func (stack *Stack) SerializeArray() []interface{} {
	arr := []interface{}{}
	next := stack.Bottom
	for next != nil {
		// fmt.Println(next.value)
		arr = append(arr, next.Value)
		next = next.Next
	}

	return arr
}

func (stack *Stack) MarshalJSON() ([]byte, error) {
	values := []interface{}{}
	currentNode := stack.Bottom
	for currentNode != nil {
		values = append(values, currentNode.Value)
		currentNode = currentNode.Next
	}
	return json.Marshal(values)
}

func (stack *Stack) UnmarshalJSON(data []byte) error {
	var values []interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	stack.Length = len(values)
	for i := range values {
		value := values[i]
		node := &StackNode{Value: value}
		if i == 0 {
			stack.Bottom = node
			stack.Top = node
		} else {
			stack.Top.Next = node
			node.Prev = stack.Top
			stack.Top = node
		}
	}

	// fmt.Println("unmarshal: ", values)
	// fmt.Println("unmarshal[0]: ", values[0])
	fmt.Println("unmarshal: ", stack.Bottom)
	return nil
}

// type ListNode struct {
// 	value interface{}
// 	next  *ListNode
// }

// type LinkedList struct {
// 	head *ListNode
// 	tail *ListNode
// }

// func (list *LinkedList) Push(value interface{}) {
// 	newList := &ListNode{value: value}
// 	if list.head == nil {
// 		list.head = newList
// 		list.tail = newList
// 	} else {
// 		list.tail.Next = newList
// 		list.tail = newList
// 	}
// }
