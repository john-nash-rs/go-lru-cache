package cache

import (
	"hash/fnv"
)

type Cache[T any] struct {
	doubleLinkedList *List[*Item[T]]
	maxSize          int32
	buckets          []*Bucket[T]
	bucketCount      int32
	currentSize      int32
}

func New[T any](listSize int32, bucketCount int32) *Cache[T] {
	c := &Cache[T]{
		doubleLinkedList: &List[*Item[T]]{},
		maxSize:          listSize,
		buckets:          make([]*Bucket[T], bucketCount),
		bucketCount:      bucketCount,
	}

	for i := int32(0); i < bucketCount; i++ {
		c.buckets[i] = &Bucket[T]{
			kvstore: make(map[string]*Item[T]),
		}
	}
	return c
}

func (c *Cache[T]) Get(key string) T {
	bucket := c.findBucket(key, c.bucketCount)
	bucket.lock.RLocker().Lock()
	defer bucket.lock.RLocker().Unlock()
	return bucket.kvstore[key].value
}

func (c *Cache[T]) Put(key string, value T) {
	//First find bucket to store the data
	bucket := c.findBucket(key, c.bucketCount)
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	item := &Item[T]{
		key:   key,
		value: value,
	}
	bucket.kvstore[key] = item
	//Promote(move) the key to head of linked list as this is the most recently used
	c.promote(item)
	c.currentSize = c.currentSize + 1
}

func (c *Cache[T]) promote(item *Item[T]) {
	if c.currentSize > c.maxSize {
		itemToDelete := c.doubleLinkedList.tail
		c.removeFromMap(itemToDelete.item.key)
		c.doubleLinkedList.tail = c.doubleLinkedList.tail.next
		node := &Node[*Item[T]]{
			item:     item,
			previous: c.doubleLinkedList.head.previous,
		}
		c.doubleLinkedList.head.next = node
		prevNode := c.doubleLinkedList.head
		c.doubleLinkedList.head = node
		c.doubleLinkedList.head.previous = prevNode

		return
	}

	if c.currentSize == 0 {
		c.doubleLinkedList.head = &Node[*Item[T]]{
			item: item,
		}
		c.doubleLinkedList.tail = c.doubleLinkedList.head
		return
	}

	node := &Node[*Item[T]]{
		item:     item,
		previous: c.doubleLinkedList.head.previous,
	}

	c.doubleLinkedList.head.next = node
	prevNode := c.doubleLinkedList.head
	c.doubleLinkedList.head = node
	c.doubleLinkedList.head.previous = prevNode

}

func (c *Cache[T]) removeFromMap(key string) {
	bucket := c.findBucket(key, c.bucketCount)
	delete(bucket.kvstore, key)

}

func (c *Cache[T]) findBucket(key string, bucketCount int32) *Bucket[T] {
	h := fnv.New32a()
	h.Write([]byte(key))
	bucketIndex := int32(h.Sum32()) % bucketCount
	return c.buckets[bucketIndex]
}
