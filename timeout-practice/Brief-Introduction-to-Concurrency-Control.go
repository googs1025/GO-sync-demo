package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
	https://mp.weixin.qq.com/s/uKT_yU7j8ILNweygFnCiSQ
	WaitGroup：多个goroutine的任务处理存在依赖或拼接关系。
	channel+select：可以主动取消goroutine；多groutine中数据传递；channel可以代替WaitGroup的工作，但会增加代码逻辑复杂性；多channel可以满足Context的功能，同样，也会让代码逻辑变得复杂。
	Context：多层级groutine之间的信号传播（包括元数据传播，取消信号传播、超时控制等）。
 */


func main() {
	//WaitgroupTry()
	//ChannelExitTry()
	ContextTry()
	
}

// 某任务需要多 goroutine 协同工作，每个 goroutine 只能做该任务的一部分，只有全部的 goroutine 都完成，任务才算是完成。
func WaitgroupTry() {
	var wg sync.WaitGroup

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(index int) {

			defer wg.Done()
			fmt.Println("任务执行中！", index)

		}(i)

	}

	// 主goroutine会在这里阻塞
	wg.Wait()

	fmt.Println("所有任务执行完毕")
}

// channel+select 的组合，是优雅的通知goroutine 结束的方式
// 缺点：如果有多个 goroutine 都需要控制结束怎么办？如果这些 goroutine 又衍生了其它更多的goroutine呢？
// 调用链复杂时，不推荐用chan+select 方式，需要用context
func ChannelExitTry() {
	stop := make(chan bool)

	go func() {

		for {
			select {
			case <-stop:
				fmt.Println("收到结束通知")
				return
			default:
				fmt.Println("持续监控中")
				time.Sleep(time.Second * 2)

			}
		}

	}()

	time.Sleep(time.Second * 5)
	fmt.Println("通知子gorotine退出")
	stop <- true

	time.Sleep(time.Second * 3)

}


func A(ctx context.Context, name string) {
	go B(ctx, name)

	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "A退出")
			return
		default:
			// 业务逻辑
			fmt.Println(name, "A do something")
			time.Sleep(time.Second * 2)
		}
	}

}

func B(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "B退出")
			return
		default:
			// 业务逻辑
			fmt.Println(name, "B do something")
			time.Sleep(time.Second * 2)
		}
	}
}

func ContextTry() {
	// 调用context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// 开始调用子goroutine
	go A(ctx, "[请求1]")
	time.Sleep(3 * time.Second)
	fmt.Println("client断开连接，通知对应处理client请求的A,B退出")
	// 调用删除方法
	cancel()
	time.Sleep(3 * time.Second)

}


