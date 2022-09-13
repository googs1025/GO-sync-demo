package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/*

 */

const (
	Max = 100000
	NumReceivers = 10
	NumSenders = 100
)

var wgReceivers sync.WaitGroup

func main() {

	rand.Seed(time.Now().UnixNano())

	// 数据缓存chan
	dataCh := make(chan int, 100)
	// 通知退出chan
	stopCh := make(chan struct{})


	// toStop := make(chan string, NumReceivers + NumReceivers)
	// toStop := make(chan string, 1)
	// 两种都可以 都不会阻塞

	toStop := make(chan string, 1)
	var stoppedBy string


	go func() {
		stoppedBy = <-toStop
		close(stopCh)

	}()

	// 生产者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				//
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}
				select {
				case <-stopCh:
					return
				case dataCh <-value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 消费者
	for i := 0; i < NumReceivers; i++ {
		// 消费者需要用waitgroup兜底
		wgReceivers.Add(1)
		go func(id string, wgReceivers *sync.WaitGroup) {
			defer wgReceivers.Done()
			for {
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:

						}
						return
					}
					fmt.Println("value:", value)
				}
			}

		}(strconv.Itoa(i), &wgReceivers)
	}

	wgReceivers.Wait()
	fmt.Println("stopped by", stoppedBy)



}
