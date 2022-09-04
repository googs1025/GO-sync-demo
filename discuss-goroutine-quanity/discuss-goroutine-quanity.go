package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

// https://mp.weixin.qq.com/s?__biz=MzUzMTUxMzYyNQ==&mid=2247483760&idx=1&sn=65632f35d4d0bf11e3f7f4c1d99e0549&chksm=fa402906cd37a010a63f0d961662513eb379d4ef00ec266bed886c0b52f31fc667e53c299ff1&cur_album_id=1870590690292793347&scene=190#rd

func main() {

	//UnlimitedNumberofGoroutines()
	//LimitedNumberofGoroutine1()
	LimitedNumberofGoroutine2()

}


// 1. 没有限制goroutine ！

func UnlimitedNumberofGoroutines() {
	task_num := math.MaxInt
	for i := 0; i < task_num; i++ {
		go func(i int) {
			fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
		}(i)
	}
}


//  2. 采用chan+waitgroup的方式 限制goroutine的数量。

var wg sync.WaitGroup

func aaa(taskChan chan bool, num int) {
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Println("go func ", num, " goroutine count = ", runtime.NumGoroutine())
	<- taskChan


}

func LimitedNumberofGoroutine1() {

	task_num := 100
	taskChan := make(chan bool, 6)


	for i := 0; i < task_num; i++ {
		wg.Add(1)
		taskChan <-true
		go aaa(taskChan, i)
	}

	wg.Wait()
	fmt.Println("主goroutine退出")

}


// 3.利用无缓冲channel与任务发送/执行分离方式

var wg2 sync.WaitGroup

func LimitedNumberofGoroutine2() {

	ch := make(chan int)
	goroutine_num := 5
	for i := 0; i < goroutine_num; i++ {
		go example(ch)
	}

	task_num := 10
	//task_num := math.MaxInt

	for t := 0; t < task_num; t++ {
		taskSend(t, ch)

	}
	wg2.Wait()

}

func taskSend(task int, ch chan int) {
	wg2.Add(1)
	ch <-task

}

// 执行任务
func example(ch chan int) {

	for t := range ch {
		// 执行业务逻辑
		fmt.Println("go task = ", t, ", goroutine count = ", runtime.NumGoroutine())
		wg2.Done()
	}

}


