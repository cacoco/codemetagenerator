package utils

type (
	Stack[T any] struct {
		top    *node[T]
		length int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

var NilStack = &Stack[string]{nil, 0}

// Create a new stack
func New[T any]() *Stack[T] {
	return &Stack[T]{nil, 0}
}

// Return the number of items in the stack
func (s *Stack[T]) Len() int {
	return s.length
}

// View the top item on the stack
func (s *Stack[T]) Peek() *T {
	if s.length == 0 {
		return nil
	}
	return &s.top.value
}

// Pop the top item of the stack and return it
func (s *Stack[T]) Pop() *T {
	if s.length == 0 {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--
	return &n.value
}

// Push a value onto the top of the stack
func (s *Stack[T]) Push(value T) {
	n := &node[T]{value, s.top}
	s.top = n
	s.length++
}
