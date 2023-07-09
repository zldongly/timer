package timer

import (
	"time"
)

type Task struct {
	id int64

	expire time.Time
	fn     func()

	next *Task // 链表
}

type chain struct {
	head *Task
}

// insert
// 按expire插入
func (c *chain) insert(t *Task) {
	if c.head == nil {
		c.head = t
		return
	}

	if t.expire.Before(c.head.expire) {
		t.next = c.head
		c.head = t
		return
	}

	var curr *Task
	for curr = c.head; curr.next != nil; curr = curr.next {
		//if !curr.next.expire.Before(t.expire) {
		if t.expire.Before(curr.next.expire) {
			t.next = curr.next
			curr.next = t
			return
		}
	}

	curr.next = t
}

// del
// 按id删除
// return被删除的那个节点
func (c *chain) del(id int64) *Task {
	if c.head == nil {
		return nil
	}

	var (
		prev *Task = nil
		curr       = c.head
	)

	for curr != nil {
		if curr.id == id {
			if prev == nil { // 删除head
				c.head = curr.next

				curr.next = nil
				return curr
			}

			prev.next = curr.next

			curr.next = nil
			return curr
		}

		prev = curr
		curr = curr.next
	}

	return nil
}

// pop
// 取出head
func (c *chain) pop() *Task {
	if c.head == nil {
		return nil
	}

	t := c.head
	c.head = t.next

	return t
}

// stack
type stack struct {
	head *Task
	size int
}

// push stack
func (s *stack) push(t *Task) {
	s.size++

	t.next = s.head
	s.head = t
}

// pop  stack
func (s *stack) pop() *Task {
	if s.head == nil {
		return nil
	}
	s.size--

	t := s.head
	s.head = s.head.next
	t.next = nil
	return t
}
