package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
	多个生产者+一个消费者的模式
	不要close 缓存的dataCh，使用stopC从消费者端通知。

	注：不要全部都用waitgroup来做，只要消费者用即可。
	因为生产者有close兜底了，所以只要确认消费者的退出就好。
 */


// 常数
const Max = 10000
const NumSenders = 100

var wgReceivers sync.WaitGroup

func main() {

	rand.Seed(time.Now().UnixNano())

	// 数据缓存
	dataC := make(chan int, 100)
	// 通知退出
	stopC := make(chan struct{})

	// 多个生产者
	for i := 0; i < NumSenders; i++ {
		go func(dataC chan int) {


			for {
				select {
				case dataC <-rand.Intn(Max):	// 生产数据
				case <-stopC:	// 通知退出
					return
				}
			}
		}(dataC)
	}

	// 消费者
	wgReceivers.Add(1)
	go func(dataC chan int, wgReceivers *sync.WaitGroup) {
		// 遍历
		defer wgReceivers.Done()

		for value := range dataC {

			// 一个退出的条件，由消费者来控制退出！
			if value == Max-1 {
				fmt.Println("send stop signal to senders.")
				close(stopC)
				return
			}
			fmt.Println("value:", value)
		}

	}(dataC, &wgReceivers)

	wgReceivers.Wait()


}
