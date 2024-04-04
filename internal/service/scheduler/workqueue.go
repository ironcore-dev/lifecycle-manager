// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"fmt"
	"sync"
)

type RingBufQueue[T LifecycleObject] struct {
	buf  []Task[T]
	cap  uint64
	head uint64
	tail uint64
	full bool

	keyToIndex map[string]uint64
	indexToKey map[uint64]string

	Enqueued chan struct{}

	mu sync.RWMutex
}

func NewRingBufQueue[T LifecycleObject](capacity uint64) *RingBufQueue[T] {
	return &RingBufQueue[T]{
		keyToIndex: make(map[string]uint64),
		indexToKey: make(map[uint64]string),
		buf:        make([]Task[T], capacity),
		cap:        capacity,
		Enqueued:   make(chan struct{}, capacity),
	}
}

func (q *RingBufQueue[T]) Enqueue(item Task[T]) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return false
	}

	q.buf[q.tail] = item
	q.keyToIndex[item.Key] = q.tail
	q.indexToKey[q.tail] = item.Key
	q.tail = (q.tail + 1) % q.cap
	q.full = q.head == q.tail
	q.Enqueued <- struct{}{}
	return true
}

func (q *RingBufQueue[T]) Dequeue() (Task[T], bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var item Task[T]
	if q.IsEmpty() {
		return item, false
	}

	item = q.buf[q.head]
	key := q.indexToKey[q.head]
	delete(q.indexToKey, q.head)
	delete(q.keyToIndex, key)
	q.head = (q.head + 1) % q.cap
	q.full = false
	return item, true
}

func (q *RingBufQueue[T]) IsEmpty() bool {
	return !q.full && q.head == q.tail
}

func (q *RingBufQueue[T]) IsFull() bool {
	return q.full
}

func (q *RingBufQueue[T]) Has(key string) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	_, stored := q.keyToIndex[key]
	return stored
}

func (q *RingBufQueue[T]) Len() int {
	return len(q.keyToIndex)
}

func (q *RingBufQueue[T]) FreeCapacity() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return int(q.cap) - len(q.keyToIndex)
}

func (q *RingBufQueue[T]) TryEnqueue(item Task[T]) bool {
	if q.IsFull() {
		return false
	}
	return q.Enqueue(item)
}

func (q *RingBufQueue[T]) Print() string {
	return fmt.Sprint(q.keyToIndex)
}
