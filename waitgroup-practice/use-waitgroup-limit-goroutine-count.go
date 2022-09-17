package waitgroup_practice

import (
	"fmt"
	"net/http"
	"sync"
	//"time"

)

/*
	重新定义WaitGroup，目的就是为了支持并发数量限制，跟以往不确定并发相比，重新构造的WaitGroup可以限制并发数量和查看pending数量。
 */

// 重新封装waitgroup
type WaitGroup struct {
	// 主结构体，复用并发控制
	waitGroup 	sync.WaitGroup
	// chan大小，用来限制gorotuine的数量
	size 		int
	// chan，用来规定一段时间只能有一定的goroutine，否则阻塞
	pool 		chan struct{}

}

// 建立新对象
func NewWaitGroup(size int) *WaitGroup {
	// 对象
	wg := &WaitGroup{
		size: size,
	}
	// size大于0，初始化chan
	if size > 0 {
		wg.pool = make(chan struct{}, size)
	}

	return wg

}

// 就是Add()方法，只是当size大于0时，需要发一个byte给wg.pool chan，相当于占用位置
func (wg *WaitGroup) BlockAdd() {
	if wg.size > 0 {
		wg.pool <-struct{}{}
	}
	wg.waitGroup.Add(1)
}

// 就是Done()方法，只是当size大于0时，需要给释放掉
func (wg *WaitGroup) Done() {
	if wg.size > 0 {
		<- wg.pool
	}
	wg.waitGroup.Done()
}

func (wg *WaitGroup) Wait() {
	wg.waitGroup.Wait()
}

func (wg *WaitGroup) PendingCount() int64 {
	return int64(len(wg.pool))
}

func WaitGroupPractice3() {
	urls := []string{
		"https://www.a.com/",
		"https://www.b.com",
		"https://www.c.com",
		"https://www.d.com/",
		"https://www.e.com",
		"https://www.f.com",
	}
	//timer := time.NewTicker(time.Second)

	// 最多只让3个goroutine执行
	wg := NewWaitGroup(3)

	for _, url := range urls {
		wg.BlockAdd()
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("%s: result: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

		}(url)
	}

	//想要做一个定时查看状态的任务。
	//for {
	//	select {
	//	case <- timer.C:
	//		fmt.Println(wg.PendingCount())
	//
	//
	//	}
	//}


	wg.Wait()

	fmt.Println("Finished")
}
