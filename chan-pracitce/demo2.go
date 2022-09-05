package main

import (
	"fmt"
	"math/rand"
)

// 建立多个Producer，并用chan传递data，再用一个函数 for-select merge起来。

func Producer1() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <-rand.Int()
		}
	}()
	return ch
}

func Producer2() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <-rand.Int()
		}
	}()
	return ch

}

func MergeProducer() chan int {

	ch := make(chan int, 30)
	go func() {
		for {
			select {
			case v1 := <-Producer1():
				ch <- v1
			case v2 := <-Producer2():
				ch <- v2
			}
		}
	}()

	return ch


}

func main() {

	ch := MergeProducer()
	for i := 0; i < 200; i++ {
		fmt.Println(i, <-ch)
	}



}
