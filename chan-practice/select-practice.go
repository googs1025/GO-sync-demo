package main

import (
	"fmt"
	"sync"

)

// https://www.cnblogs.com/f-ck-need-u/p/9986335.html

func main() {
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	stopC := make(chan struct{})
    var wg sync.WaitGroup

	wg.Add(2)
	go pump1(ch1, &wg)
	go pump2(ch2, &wg)
	go receive(ch1, ch2, stopC)

	wg.Wait()
	stopC <- struct{}{}


}

func pump1(ch chan int, wg *sync.WaitGroup) {
	for i := 0; i <= 30; i++ {
		if i%2 == 0 {
			ch <-i
		}
	}
	wg.Done()
	close(ch)
}

func pump2(ch chan int, wg *sync.WaitGroup) {
	for i := 0; i <= 30; i++ {
		if i%2 == 1 {
			ch <-i
		}
	}
	wg.Done()
	close(ch)

}

/*
	如果在select中执行send操作，则可能会永远被send阻塞。
	所以，在使用send的时候，应该也使用defalut语句块，保证send不会被阻塞。
	如果没有default，或者能确保select不阻塞的语句块，则迟早会被send阻塞。
 */
func receive(ch1, ch2 chan int, stopC chan struct{}) {

	for {
		select {
		case v, ok := <-ch1:
			if ok {
				fmt.Printf("Recv on ch1: %d\n", v)
			}
		case v, ok := <-ch2:
			if ok {
				fmt.Printf("Recv on ch2: %d\n", v)
			}
		case <-stopC:
			fmt.Println("消费者消费完毕")
			return
		default:
			fmt.Println("waiting")
		}
	}


}
