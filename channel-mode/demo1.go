package main

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


func main() {

	ch := make(chan int, 50)
	var wg sync.WaitGroup

	wg.Add(3)
	go Producer(ch, &wg)
	go Consumer(ch, &wg)
	go Consumer(ch, &wg)

	stopC := make(chan os.Signal, 1)


	wg.Wait()

	signal.Notify(stopC, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("quit (%v)\n", <-stopC)
	fmt.Println("主goroutine退出")
}
