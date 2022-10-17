package waitgroup_practice

import (
	"fmt"
	"sync"
	"time"
)


// 业务函数
func job(index int) int {
	// 模拟耗时操作
	time.Sleep(time.Millisecond * 500)
	return index
}

func NoGoroutine() {
	start := time.Now()
	num := 5
	for i := 0; i < num ; i++ {
		fmt.Println(job(i))
	}

	end := time.Since(start)
	fmt.Println("耗时:", end.String())
}


func WaitGroupPractice4() {
	wg := sync.WaitGroup{}
	start := time.Now()
	num := 5

	for i := 0; i < num ; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Println(job(index))
		}(i)

	}

	wg.Wait()
	end := time.Since(start)
	fmt.Println("耗时:", end.String())
}

func WaitGroupPractice5() {
	wg := sync.WaitGroup{}
	start := time.Now()
	/*
		这里有个坑：如果chan是make(chan int) 同步，就不需要WaitGroup(会死锁)，直接后面取chan数据时，兜底判断if count == num 即可
		如果使用make(chan int, 5) 可以搭配WaitGroup
	 */
	num := 5
	resC := make(chan int, 5)
	for i := 0; i < num ; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			resC <-job(index)
		}(i)

	}

	wg.Wait()

	end := time.Since(start)
	fmt.Println("耗时:", end.String())

	count := 0
	for i := range resC {
		count++
		fmt.Println("收到chan中的结果:", i)
		if count == num {
			close(resC)	// 关闭chan
			break
		}
	}



}


func WaitGroupPractice6() {
	wg := sync.WaitGroup{}
	start := time.Now()

	num := 5
	resC := make(chan int)
	for i := 0; i < num ; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			resC <-job(index)
		}(i)

	}

	// 记得可以这样写。！！
	go func() {
		defer close(resC)
		wg.Wait()	// wg.Wait() 可以写在子goroutine中。
	}()


	end := time.Since(start)
	fmt.Println("耗时:", end.String())

	count := 0
	for i := range resC {
		count++
		fmt.Println("收到chan中的结果:", i)

	}



}
