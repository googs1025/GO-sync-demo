package channel_mode

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
	生产者消费者模式+waitgroup+优雅退出。
 */

func Producer(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		value := i
		if value % 2 == 0 {
			ch <-value
		}
	}
	close(ch)


}

func Consumer(inputC chan int, wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		v, ok := <-inputC
		time.Sleep(time.Second)
		if !ok {
			break
		}
		fmt.Println(v)
	}

	return
}


func ChanDemo1() {

	ch := make(chan int, 50)
	var wg sync.WaitGroup
	//ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	wg.Add(3)
	go Producer(ch, &wg)
	go Consumer(ch, &wg)
	go Consumer(ch, &wg)



	stopC := make(chan os.Signal)
	signal.Notify(stopC, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopC:
		return
	default:
		wg.Wait()
	}

	fmt.Println("主进程完成生产者消费者任务")

}
