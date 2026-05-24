package lfucache

import (
	"errors"
	"iter"

	ll "github.com/igoroutine-courses/gonature.lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Put(key K, value V)
	All() iter.Seq2[K, V]
	Size() int
	Capacity() int
	GetKeyFrequency(key K) (int, error)
}

var _ Cache[int, int] = (*cacheImpl[int, int])(nil)

type cacheImpl[K comparable, V any] struct {
	linkedList   ll.List[node[K, V]]
	keyToElement map[K]*ll.Element[node[K, V]]
	lastUpdated  map[int]*ll.Element[node[K, V]]
	capacity     int
}

type node[K comparable, V any] struct {
	key   K
	value V
	count int
}

func (l *cacheImpl[K, V]) incrementFrequency(element *ll.Element[node[K, V]]) {
	count := element.Value.count
	if l.lastUpdated[count] == element {
		if element.Next() != nil && element.Next().Value.count == count {
			l.lastUpdated[count] = element.Next()
		} else {
			delete(l.lastUpdated, count)
		}
	}

	if _, ok := l.lastUpdated[count]; ok {
		l.linkedList.MoveBefore(element, l.lastUpdated[count])
	}

	if el, ok := l.lastUpdated[count+1]; ok {
		l.linkedList.MoveBefore(element, el)
	}

	l.lastUpdated[count+1] = element
	element.Value.count++
}

func New[K comparable, V any](capacity ...int) *cacheImpl[K, V] {
	realCapacity := DefaultCapacity
	if len(capacity) > 0 {
		realCapacity = capacity[0]
	}
	if realCapacity < 0 {
		panic("Negative capacity")
	}
	return &cacheImpl[K, V]{
		linkedList:   ll.NewList[node[K, V]](),
		keyToElement: make(map[K]*ll.Element[node[K, V]], realCapacity),
		lastUpdated:  make(map[int]*ll.Element[node[K, V]]),
		capacity:     realCapacity,
	}
}

func (l *cacheImpl[K, V]) Get(key K) (V, error) {
	if el, ok := l.keyToElement[key]; ok {
		l.incrementFrequency(el)
		return el.Value.value, nil
	}
	return *new(V), ErrKeyNotFound
}

func (l *cacheImpl[K, V]) Put(key K, value V) {
	if el, ok := l.keyToElement[key]; ok {
		el.Value.value = value
		l.incrementFrequency(el)
		return
	}

	n := node[K, V]{
		key:   key,
		value: value,
		count: 0,
	}

	var el *ll.Element[node[K, V]]
	if l.Size() == l.Capacity() {
		el = l.linkedList.Back()

		if l.lastUpdated[el.Value.count] == el {
			delete(l.lastUpdated, el.Value.count)
		}
		delete(l.keyToElement, el.Value.key)
	} else {
		el = l.linkedList.PushBack(n)
	}
	el.Value = n
	l.keyToElement[key] = el
	l.incrementFrequency(l.keyToElement[key])
}

func (l *cacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		cur := l.linkedList.Front()

		for range l.Size() {
			n := cur.Value

			if !yield(n.key, n.value) {
				return
			}
			cur = cur.Next()
		}
	}
}

func (l *cacheImpl[K, V]) Size() int {
	return len(l.keyToElement)
}

func (l *cacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *cacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	if el, ok := l.keyToElement[key]; ok {
		return el.Value.count, nil
	}
	return 0, ErrKeyNotFound
}
