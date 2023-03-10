package ds

import (
	"encoding/json"
	"fmt"
	"log"
)

type Node struct {
	Value interface{}
	Next  *Node
	Prev  *Node
}

// Doubled Linked List
type LinkedList struct {
	Head *Node
	Tail *Node
}

func (list *LinkedList) Push(value interface{}) {
	newList := &Node{Value: value}
	if list.Head == nil {
		list.Head = newList
		list.Tail = newList
	} else {
		newList.Prev = list.Tail
		list.Tail.Next = newList
		list.Tail = newList
	}
}

func (list *LinkedList) Remove(node *Node) {
	if list.Head == node {
		list.Head = node.Next
		list.Head.Prev = nil
	} else if list.Tail == node {
		list.Tail = node.Prev
		list.Tail.Next = nil
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
}

func (list *LinkedList) Print() {
	currentNode := list.Head
	for currentNode != nil {
		log.Println(currentNode.Value)
		currentNode = currentNode.Next
	}
}

// Stack using Doubly Linked List
type Stack struct {
	Bottom *Node
	Top    *Node
	Length int
}

func (stack *Stack) Push(value interface{}) {
	newStack := &Node{Value: value}
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
		node := &Node{Value: value}
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

func (stack *Stack) Print() {
	currentNode := stack.Bottom
	for currentNode != nil {
		log.Println(currentNode.Value)
		currentNode = currentNode.Next
	}
}
