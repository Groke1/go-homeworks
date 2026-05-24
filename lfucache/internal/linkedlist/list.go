package linkedlist

type List[T any] interface {
	Front() *Element[T]
	Back() *Element[T]
	PushBack(value T) *Element[T]
	MoveBefore(element1, element2 *Element[T])
	Remove(*Element[T])
	IsEmpty() bool
}

type Element[T any] struct {
	prev, next *Element[T]
	Value      T
}

func (e *Element[T]) Next() *Element[T] {
	return e.next
}

func (e *Element[T]) Prev() *Element[T] {
	return e.prev
}

func newElement[T any](value T) *Element[T] {
	return &Element[T]{prev: nil, next: nil, Value: value}
}

func (e *Element[T]) Remove() {
	if e.Prev() != nil {
		e.Prev().next = e.next
	}
	if e.Next() != nil {
		e.Next().prev = e.prev
	}
	e.prev = nil
	e.next = nil
}

type linkedList[T any] struct {
	root *Element[T]
}

func NewList[T any]() List[T] {
	return &linkedList[T]{}
}

func (l *linkedList[T]) IsEmpty() bool {
	return l.root == nil
}

func (l *linkedList[T]) Front() *Element[T] {
	if l.IsEmpty() {
		return nil
	}
	return l.root.Next()
}

func (l *linkedList[T]) Back() *Element[T] {
	if l.IsEmpty() {
		return nil
	}
	return l.root.Prev()
}

func (l *linkedList[T]) mergeElements(element1, element2 *Element[T]) {
	element1.next, element2.prev = element2, element1
}

func (l *linkedList[T]) PushBack(value T) *Element[T] {
	el := newElement[T](value)
	if l.IsEmpty() {
		l.root = &Element[T]{prev: nil, next: nil}
		l.mergeElements(l.root, el)
		l.mergeElements(el, l.root)
	} else {
		l.MoveBefore(el, l.root)
	}
	return el
}

func (l *linkedList[T]) MoveBefore(element1, element2 *Element[T]) {
	element1.Remove()
	l.mergeElements(element2.prev, element1)
	l.mergeElements(element1, element2)
}

func (l *linkedList[T]) Remove(element *Element[T]) {
	element.Remove()
}
