# Timer 
用于执行定时任务  
需要执行n个定时任务时，根据任务启动时间排序，多个定时任务共用同一个定时器，不需要每个任务启一个定时器，减少goroutine资源  

## Example
```go
package main

import (
	"github.com/zldongly/timer"
	"log"
	"time"
)

func main() {
	tr := timer.New()

	// 添加一个定时任务
	id := tr.Add(time.Now().Add(time.Second*10), func() {
		// TODO
	})

	// 重新定时
	err := tr.Set(id, time.Now().Add(time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	// 删除定时任务
	tr.Del(id)
    
	// 循环启动定时任务
	loop(tr)
}

func loop(tr *timer.Timer) {
	tr.Add(time.Now().Add(time.Second), func() {
		// TODO

		loop(tr)
	})
}
```
