package pool_practice

import (
	"fmt"
	"time"
)

// https://mp.weixin.qq.com/s/7gjU68kajH1pDZLw9WJ5zA

/*
1、如果设定的协程池数大于 2，
此时第二次传入往 NewTask 传入task，select case 的时候，
如果第一个协程还在运行中，就一定会走第二个case，重新创建一个协程执行task

2、如果传入的任务数大于设定的协程池数，并且此时所有的任务都还在运行中，
那此时再调用 NewTask 传入 task ，这两个 case 都不会命中，
会一直阻塞直到有任务执行完成，worker 函数里的 work 通道才能接收到新的任务，继续执行。
 */

// 协程池
type Pool struct {
	// 一个是 work，用于接收 task 任务
	//一个是 sem，用于设置协程池大小，即可同时执行的协程数量
	work chan func()
	sem chan struct{}
}

// 创建pool对象
func NewPool(size int) *Pool {
	// worker 是一个无缓冲chan
	// sem 是一个有size 缓冲的chan，size大小即是池大小
	return &Pool{
		work: make(chan func()),
		sem: make(chan struct{}, size),
	}
}
// 启新任务
// 当第一次调用 NewTask 添加任务的时候，由于 work 是无缓冲通道，
// 所以一定会走第二个 case 的分支：使用 go worker 开启一个协程。
func (p *Pool) NewTask(task func()) {
	select {
	case p.work <-task: // 第一次调用一定会把task放进chan中
	case p.sem <- struct{}{}:   // 当p.work中已经阻塞时，就会调用这里，启goroutine来跑了
		go p.worker(task)

	}
}
// 为了能够实现协程的复用，这个使用了 for 无限循环，
// 使这个协程在执行完任务后，也不退出，而是一直在接收新的任务。
// 这里要注意一下！！
func (p *Pool) worker(task func()) {
	//执行完
	defer func() {<- p.sem}()
	for {
		// 执行任务

		task()	// 启后需要先执行，再次把
		task = <-p.work

	}
}


func PoolPractice1() {

	//pool := NewPool(128)
	//pool.NewTask(func() {
	//	fmt.Println("doing task!!")
	//})
	//
	//time.Sleep(5 * time.Second)

	// 启动五个goroutine，以及一个chan大小是5的池
	pool := NewPool(5)

	for i := 1; i <50; i++{

		pool.NewTask(func(){
			time.Sleep(2 * time.Second)
			fmt.Println(time.Now())
		})
	}


	// 保证所有的协程都执行完毕
	time.Sleep(5 * time.Second)

	
}
