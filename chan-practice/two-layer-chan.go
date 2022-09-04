package main

import (
	"fmt"
	"time"
)

// https://www.cnblogs.com/f-ck-need-u/p/9994496.html

func main() {
	cc := make(chan chan int)

	times := 5
	for i := 1; i < times; i++ {
		f := make(chan bool)

		go f1(cc, f)

		ch := <-cc
		ch <-i

		for sum := range ch {
			fmt.Printf("Sum(%d)=%d\n", i, sum)
		}

		time.Sleep(time.Second)
		close(f)
	}
}

func f1(cc chan chan int, f chan bool) {
	c := make(chan int)
	cc <-c
	defer close(c)
	sum := 0

	select {
	case x := <-c:
		for i := 0; i <=x; i++ {
			sum = sum + i
		}
		c <-sum
	case <-f:
		return
	}
}