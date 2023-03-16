package cache

import "sync"

type List[T any] struct {
	head *Node[T]
	tail *Node[T]
}

type Node[T any] struct {
	next     *Node[T]
	previous *Node[T]
	item     T
}

type Item[T any] struct {
	key   string
	value T
}

type Bucket[T any] struct {
	lock    sync.RWMutex
	kvstore map[string]*Item[T]
}
