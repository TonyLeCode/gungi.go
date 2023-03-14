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

func (l *LinkedList) Push(value interface{}) *Node {
	newNode := &Node{Value: value}
	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
	} else {
		newNode.Prev = l.Tail
		l.Tail.Next = newNode
		l.Tail = newNode
	}
	return newNode
}

func (l *LinkedList) Remove(node *Node) {
	if l.Head == node {
		l.Head = node.Next
		l.Head.Prev = nil
	} else if l.Tail == node {
		l.Tail = node.Prev
		l.Tail.Next = nil
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
}

func (l *LinkedList) Print() {
	currentNode := l.Head
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

func (s *Stack) Push(value interface{}) {
	newNode := &Node{Value: value}
	if s.Bottom == nil {
		s.Bottom = newNode
		s.Top = newNode
	} else {
		s.Top.Next = newNode
		newNode.Prev = s.Top
		s.Top = newNode
	}
	s.Length++
}

func (s *Stack) Pop() {
	if s.Top == nil {
		return
	}
	s.Top = s.Top.Prev
	if s.Top != nil {
		s.Top.Next = nil
	} else {
		s.Bottom = nil
	}
	s.Length--
}

func (s *Stack) SerializeArray() []interface{} {
	arr := []interface{}{}
	next := s.Bottom
	for next != nil {
		// fmt.Println(next.value)
		arr = append(arr, next.Value)
		next = next.Next
	}

	return arr
}

func (s *Stack) MarshalJSON() ([]byte, error) {
	values := []interface{}{}
	currentNode := s.Bottom
	for currentNode != nil {
		values = append(values, currentNode.Value)
		currentNode = currentNode.Next
	}
	return json.Marshal(values)
}

func (s *Stack) UnmarshalJSON(data []byte) error {
	var values []interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	s.Length = len(values)
	for i := range values {
		value := values[i]
		node := &Node{Value: value}
		if i == 0 {
			s.Bottom = node
			s.Top = node
		} else {
			s.Top.Next = node
			node.Prev = s.Top
			s.Top = node
		}
	}
	// fmt.Println("unmarshal: ", values)
	// fmt.Println("unmarshal[0]: ", values[0])
	fmt.Println("unmarshal: ", s.Bottom)
	return nil
}

func (s *Stack) Print() {
	currentNode := s.Bottom
	for currentNode != nil {
		log.Println(currentNode.Value)
		currentNode = currentNode.Next
	}
}
