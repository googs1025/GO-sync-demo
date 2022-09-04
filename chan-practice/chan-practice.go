package main

import (
	"fmt"
	"sync"
)

func main() {

	valueC := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go sender(valueC, &wg)
	go receiver(valueC, &wg)


	wg.Wait()


}

func sender(ch chan<- string, wg *sync.WaitGroup) {
	ch <- "aaaa"
	ch <- "bbbb"
	ch <- "cccc"
	ch <- "dddd"
	close(ch)
	wg.Done()
}

func receiver(ch <-chan string, wg *sync.WaitGroup) {

	for {
		res, ok := <-ch
		if !ok {
			wg.Done()
			return
		}

		fmt.Println(res)
	}


}
