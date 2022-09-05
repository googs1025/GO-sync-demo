package main

import (
	"fmt"
	"math/rand"
)

// 建立多个Producer，并用chan传递data，再用一个函数 for-select merge起来。
// 加上退出通知机制。

func Producer1(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("生产者1收到退出通知！")
				return
			default:
				ch <-rand.Int()

			}

		}
	}()
	return ch
}

func Producer2(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("生产者2收到退出通知！")
				return
			default:
				ch <-rand.Int()

			}

		}
	}()
	return ch

}

func MergeProducer(done chan struct{}) chan int {

	ch := make(chan int, 50)
	stopC := make(chan struct{})
	go func() {
		for {
			select {
			case v1 := <-Producer1(stopC):
				ch <- v1
			case v2 := <-Producer2(stopC):
				ch <- v2
			case <-done:
				close(stopC)
				fmt.Println("通知多个生产者们退出")
				fmt.Println("MergeProducer自己也退出！")
				return
			}
		}
	}()

	return ch


}

func main() {
	doneC := make(chan struct{})
	ch := MergeProducer(doneC)
	for i := 0; i < 200; i++ {
		fmt.Println(i, <-ch)
	}
	doneC <- struct{}{}
	fmt.Println("主goroutine退出")



}

