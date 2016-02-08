package util

type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value interface{}
	next  *Element
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}

	return nil
}

func (s *Stack) Peek(amount int) interface{} {
	if s.Len() > 0 {
		var current *Element

		for i := 0; i < amount+1; i++ {
			if i == 0 {
				current = s.top
			} else {
				current = current.next
			}
		}

		return current.value
	}

	return nil
}
