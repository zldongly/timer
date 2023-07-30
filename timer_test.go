package timer_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zldongly/timer"
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	var (
		tr    = timer.New(timer.Size(8))
		wants = []int{2, 3, 4, 6, 7, 9, 5}
		list  = make([]int, 0, 10)
		wg    sync.WaitGroup
	)

	wg.Add(7)
	for i := 0; i < 10; i++ {
		i := i
		e := time.Now().Add(time.Millisecond * time.Duration(i))
		tr.Add(e, func() {
			list = append(list, i+1)
			wg.Done()
		})
	}

	tr.Del(1)
	tr.Del(10)
	tr.Del(8)
	// 重置时间 不存在的任务报错
	err := tr.Set(8, time.Now().Add(time.Second))
	assert.Equal(t, timer.ErrNotExist, err)
	err = tr.Set(5, time.Now().Add(time.Millisecond*20))
	assert.Equal(t, error(nil), err)

	wg.Wait()
	for i := range wants {
		assert.Equal(t, wants[i], list[i])
	}
}
