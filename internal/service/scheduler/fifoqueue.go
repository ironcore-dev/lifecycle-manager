// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"fmt"
	"sync"
)

type Node[T LifecycleObject] struct {
	value Task[T]
	next  *Node[T]
	prev  *Node[T]
}

type FIFOQueue[T LifecycleObject] struct {
	head  *Node[T]
	tail  *Node[T]
	cap   uint64
	len   uint64
	nodes map[string]struct{}

	mu sync.RWMutex
}

func NewFIFOQueue[T LifecycleObject](capacity uint64) *FIFOQueue[T] {
	return &FIFOQueue[T]{
		cap:   capacity,
		nodes: make(map[string]struct{}),
	}
}

func (q *FIFOQueue[T]) Push(item Task[T]) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return false
	}

	node := &Node[T]{value: item}
	if q.len == 0 {
		q.head = node
		q.tail = node
		q.len++
		q.nodes[node.value.Key] = struct{}{}
		return true
	}

	tail := q.tail
	node.prev = tail
	tail.next = node
	q.tail = node
	q.len++
	q.nodes[node.value.Key] = struct{}{}

	return true
}

func (q *FIFOQueue[T]) Pop() (Task[T], bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var item Task[T]
	switch q.len {
	case 0:
		return item, false
	case 1:
		item = q.head.value
		q.head = nil
		q.tail = nil
		q.len--
		delete(q.nodes, item.Key)
		return item, true
	default:
		item = q.head.value
		q.head = q.head.next
		q.head.prev = nil
		q.len--
		delete(q.nodes, item.Key)
		return item, true
	}
}

func (q *FIFOQueue[T]) IsEmpty() bool {
	return q.len == 0
}

func (q *FIFOQueue[T]) Has(key string) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	_, stored := q.nodes[key]
	return stored
}

func (q *FIFOQueue[T]) IsFull() bool {
	return q.len == q.cap
}

func (q *FIFOQueue[T]) Len() int {
	return int(q.len)
}

func (q *FIFOQueue[T]) TryPush(item Task[T]) bool {
	if q.IsFull() {
		return false
	}
	return q.Push(item)
}

func (q *FIFOQueue[T]) Print() string {
	return fmt.Sprint(q.nodes)
}
