package utils

import (
	"fmt"
)

type Node[T any] struct {
	Next     *Node[T]
	Data     *T
	sentinel bool
}

func NewNode[T any](data *T) *Node[T] {
	return &Node[T]{
		Next:     nil,
		Data:     data,
		sentinel: false,
	}
}

func (n Node[T]) IsSentinel() bool {
	return n.sentinel
}

type LinkList[T any] struct {
	Head *Node[T]
	Tail *Node[T]
}

func NewLinkList[T any]() *LinkList[T] {
	return &LinkList[T]{Head: &Node[T]{sentinel: true}}
}

type LinkListCmp[T any] func(T, T) bool

func (l *LinkList[T]) InsertBefore(data *T, cmp LinkListCmp[*T]) {

	prev := l.Head
	runner := l.Head.Next

	for runner != nil {
		if cmp(runner.Data, data) {
			tmp := NewNode[T](data)
			prev.Next = tmp
			tmp.Next = runner
			return
		}
		prev = runner
		runner = runner.Next
	}

	prev.Next = NewNode[T](data)
}

func (list LinkList[T]) Traverse() {
	current := list.Head.Next
	for current != nil {
		fmt.Printf("%v ->\n", *current.Data)
		current = current.Next
	}
	fmt.Println("nil")
}

func (list LinkList[T]) ToList() []T {
	l := make([]T, 0)
	current := list.Head.Next

	for current != nil {
		l = append(l, *current.Data)
		current = current.Next
	}

	res := make([]T, 0, len(l))

	for i := len(l) - 1; i >= 0; i-- {
		res = append(res, l[i])
	}

	return res
}
