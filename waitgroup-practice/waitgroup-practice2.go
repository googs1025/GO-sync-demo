package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// go开发当中用到了并发协程多任务，同时收集返回多任务结果，go 协程没有直接返回，只能通过chan返回收集


func httpGet(url string, response chan string, limiter chan bool, wg *sync.WaitGroup) {
	//计数器-1
	defer wg.Done()
	//coding do business
	time.Sleep(1 * time.Second)
	//结果数据传入管道
	response <- fmt.Sprintf("http get:%s", url)
	//释放一个并发
	<-limiter
}

func collect(urls []string) []string {
	var result []string
	//执行的 这里要注意  需要指针类型传入  否则会异常
	wg := &sync.WaitGroup{}
	//并发控制
	limiter := make(chan bool, 10)
	defer close(limiter)

	response := make(chan string, 20)
	wgResponse := &sync.WaitGroup{}
	//处理结果 接收结果
	go func() {
		wgResponse.Add(1)
		for rc := range response {
			result = append(result, rc)
		}
		wgResponse.Done()
	}()
	//开启协程处理请求
	for _, url := range urls {
		//计数器
		wg.Add(1)
		//并发控制 10
		limiter <- true
		go httpGet(url, response, limiter, wg)
	}
	//发送任务
	wg.Wait()
	close(response) //关闭 并不影响接收遍历
	//处理接收结果
	wgResponse.Wait()
	return result

}


func main() {
	var urls []string
	for i := 0; i < 100; i++ {
		url := "url: " + strconv.Itoa(i)
		urls = append(urls, url)
	}

	fmt.Println(time.Now())
	result := collect(urls)
	fmt.Println(time.Now())
	fmt.Println(result)
}

