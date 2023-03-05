package util

type Node struct {
	Pre  *Node       `json:"pre"`
	Next *Node       `json:"next"`
	Data interface{} `json:"data"`
}

type List struct {
	HeadNode *Node `json:"head_node"`
}

func NewList() *List {
	return &List{}
}

func (l *List) IsEmpty() bool {
	return l.HeadNode == nil
}
func (l *List) Length() int64 {
	node := l.HeadNode
	var count int64
	for node == nil {
		node = node.Next
		count++
	}
	return count
}

func (l *List) AddHeader(data interface{}) {
	node := &Node{Data: data}
	node.Next = l.HeadNode
	l.HeadNode = node
}

func (l *List) AddTail(data interface{}) {
	tail := l.HeadNode
	for tail.Next != nil {
		tail = tail.Next
	}
	tail.Next = &Node{
		Data: data,
		Next: nil,
	}
}
