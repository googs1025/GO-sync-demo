package main

import (
	"fmt"
	"math/rand"
)

// 实现带有缓冲的生产者。

func Producer() chan int {

	ch := make(chan int, 10)

	go func() {
		for {
			ch <-rand.Int()
		}
	}()

	return ch

}

func main() {

	ch := Producer()
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}

}
