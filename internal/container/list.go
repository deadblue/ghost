package container

type node[V any] struct {
	value V
	prev  *node[V]
	next  *node[V]
}

type List[V any] struct {
	head   *node[V]
	tail   *node[V]
	cursor *node[V]
	index  int
	length int
}

// Len returns the number of items in this list.
func (l *List[V]) Len() int {
	return l.length
}

// Append adds an item at the end of the list, and set its value to `value`.
func (l *List[V]) Append(value V) *List[V] {
	if l.tail == nil {
		l.tail = &node[V]{
			value: value,
		}
		l.head = l.tail
		l.cursor = l.head
	} else {
		l.tail.next = &node[V]{
			value: value,
		}
		l.tail.next.prev = l.tail
		l.tail = l.tail.next
	}
	l.length += 1
	return l
}

func (l *List[V]) AtFirst() bool {
	return l.cursor == l.head
}

func (l *List[V]) AtLast() bool {
	return l.cursor == l.tail
}

func (l *List[V]) GoFirst() bool {
	l.cursor, l.index = l.head, 0
	return l.cursor != nil
}

func (l *List[V]) GoLast() bool {
	l.cursor, l.index = l.tail, l.length-1
	return l.cursor != nil
}

func (l *List[V]) Forward() bool {
	if l.cursor.next == nil {
		return false
	}
	l.cursor = l.cursor.next
	l.index += 1
	return true
}

func (l *List[V]) Backward() bool {
	if l.cursor.prev == nil {
		return false
	}
	l.cursor = l.cursor.prev
	l.index -= 1
	return true
}

func (l *List[V]) Get() (int, V) {
	return l.index, l.cursor.value
}
