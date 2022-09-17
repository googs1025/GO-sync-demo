package waitgroup_practice

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var responseChan = make(chan string, 15)

func httpGet(url string, limiter chan bool, wg1 *sync.WaitGroup) {

	defer wg1.Done()

	time.Sleep(time.Second * 3)
	responseChan <- fmt.Sprintf("Hello Go %s", url)
	<-limiter

}

func response() {
	for resp := range responseChan {
		fmt.Println("response:", resp)
	}
}


func WaitGroupPractice1() {

	start := time.Now()
	fmt.Println("start:", start)
	go response()

	var wg1 sync.WaitGroup
	limiter := make(chan bool, 10)

	for i := 0; i < 100; i++ {
		wg1.Add(1)
		limiter <-true
		url := "url:" + strconv.Itoa(i)
		go httpGet(url, limiter, &wg1)
	}

	wg1.Wait()

	fmt.Println("执行完毕，主goroutine退出")


}
