package timer

import (
	"sync"
	"time"
)

const (
	_defaultSize = 32
)

func New(opts ...Option) *Timer {
	tr := &Timer{
		size:  _defaultSize,
		tasks: &chain{},
		free:  &stack{},
	}

	for _, o := range opts {
		o(tr)
	}

	tr.init()

	return tr
}

type Timer struct {
	signal *time.Timer // 定时器
	mu     sync.Mutex  // 互斥锁
	size   int
	pk     int64 // 自增主键，已有mu锁

	tasks *chain // 待执行任务
	free  *stack // 空闲节点
}

// public

// Add 添加任务

// Set 重置任务时间

// Del 按id删除任务

// private

// init
// 初始化
func (tr *Timer) init() {
	tr.grow()
}

// grow 扩容
func (tr *Timer) grow() {
	var tasks = make([]Task, tr.size)

	for i, _ := range tasks {
		t := &(tasks[i])
		tr.putFree(t)
	}
}

// getFree 取一个空闲任务
func (tr *Timer) getFree() *Task {
	if tr.free.size == 0 {
		tr.grow()
	}

	return tr.free.pop()
}

// putFree 放回空闲堆栈
func (tr *Timer) putFree(t *Task) {
	if tr.free.size > tr.size*2 { // 空闲任务太多时丢弃
		return
	}
	tr.free.push(t)
}

// Option
// 可选参数
type Option func(t *Timer)

func Size(size int) Option {
	if size < 1 {
		panic("size must be greater than 0")
	}
	return func(t *Timer) {
		t.size = size
	}
}
