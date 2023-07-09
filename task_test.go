package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStack(t *testing.T) {
	var s stack
	var list = []int64{1, 82, 14, 6, 0, 23, 46}
	for _, i := range list {
		s.push(&Task{id: i})
	}

	for i := len(list) - 1; i >= 0; i-- {
		ta := s.pop()
		assert.Equal(t, list[i], ta.id)
	}
}

func TestChain(t *testing.T) {
	var (
		list = []int{3, 0, -4, 5, -3, 5}
		now  = time.Now()
		ch   = &chain{}
	)

	for i, n := range list {
		ta := &Task{
			id:     int64(i + 1),
			expire: now.Add(time.Duration(n) * time.Hour),
		}
		ch.insert(ta)
	}

	var (
		ta    *Task
		delId = []int64{3, 1, 6}
		popId = []int64{5, 2, 4}
	)

	// del
	for _, id := range delId {
		ta = ch.del(id)
		assert.Equal(t, id, ta.id)
	}

	// pop
	for _, id := range popId {
		ta = ch.pop()
		assert.Equal(t, id, ta.id)
	}

	var taskNil *Task
	assert.Equal(t, taskNil, ch.pop())
}
