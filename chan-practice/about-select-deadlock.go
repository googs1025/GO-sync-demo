package chanpractice

import (
	"fmt"
	"sync"

)

// https://mp.weixin.qq.com/s/F1RGLrh371l_NpeC42FRKw

func aboutselectdeadlock() {

	var wgd sync.WaitGroup

	foo := make(chan int, 1)
	bar := make(chan int, 1)
	// 如果是chan struct{}类型的 有"广播"功能 可以直接close() 通知所有的下游goroutine
	closed := make(chan struct{})

	wgd.Add(1)

	go func(wgd *sync.WaitGroup) {
		defer wgd.Done()
		for {
			select {
			case v := <-bar:
				foo <-v
				fmt.Println("foo channel 收到 bar channel的数据！")

			case <-closed:

				fmt.Println("收到退出通知，进程退出")
				return

			}
		}

	}(&wgd)

	bar <- 1222

	close(closed)

	wgd.Wait()

	fmt.Println(<-foo)



}
