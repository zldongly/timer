package timer

import (
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrNotExist = errors.New("timer: not exist task id")
)

const (
	_defaultSize               = 32
	_maxDuration time.Duration = math.MaxInt64 // 1<<63 - 1
)

func New(opts ...Option) *Timer {
	tr := &Timer{
		signal: time.NewTimer(_maxDuration),
		size:   _defaultSize,
		tasks:  &chain{},
		free:   &stack{},
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
func (tr *Timer) Add(expire time.Time, fn func()) int64 {
	tr.mu.Lock()
	id := tr.add(expire, fn)
	tr.mu.Unlock()

	return id
}

// Set 重置任务时间
func (tr *Timer) Set(id int64, expire time.Time) error {
	tr.mu.Lock()
	err := tr.set(id, expire)
	tr.mu.Unlock()
	return err
}

// Del 按id删除任务
func (tr *Timer) Del(id int64) {
	tr.mu.Lock()
	tr.del(id)
	tr.mu.Unlock()
}

// private

// init
// 初始化
func (tr *Timer) init() {
	tr.grow()
	go tr.run()
}

// grow 扩容
func (tr *Timer) grow() {
	var tasks = make([]Task, tr.size)

	for i, _ := range tasks {
		t := &(tasks[i])
		tr.putFree(t)
	}
}

// reload 重新计时
func (tr *Timer) reload() {
	if tr.tasks.head == nil {
		tr.signal.Reset(_maxDuration)
		return
	}

	d := time.Until(tr.tasks.head.expire)
	tr.signal.Reset(d)
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

// del 按id删除任务
func (tr *Timer) del(id int64) {
	var t *Task
	t = tr.tasks.del(id)
	if t == nil {
		return
	}

	tr.putFree(t)
	tr.reload()
}

// set 修改时间
func (tr *Timer) set(id int64, expire time.Time) error {
	var t *Task
	t = tr.tasks.del(id)
	if t == nil {
		return ErrNotExist
	}

	t.expire = expire
	tr.tasks.insert(t)
	tr.reload()
	return nil
}

// add 添加时间
func (tr *Timer) add(expire time.Time, fn func()) int64 {
	tr.pk = tr.pk + 1
	t := tr.getFree()

	t.id = tr.pk
	t.expire = expire
	t.fn = fn

	tr.tasks.insert(t)
	tr.reload()

	return t.id
}

// exec
func (tr *Timer) exec() {
	tr.mu.Lock()

	t := tr.tasks.head

	if t == nil { // 任务链表为空
		tr.mu.Unlock()
		return
	}

	if time.Now().Before(t.expire) { // 首个任务时间还没到
		tr.reload()
		tr.mu.Unlock()
		return
	}

	t = tr.tasks.pop() // 取出任务
	fn := t.fn
	tr.putFree(t) // 放回空闲链表

	tr.reload()
	tr.mu.Unlock()

	if fn != nil {
		fn()
	}
}

// run
func (tr *Timer) run() {
	for {
		<-tr.signal.C
		tr.exec()
	}
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
