package timeout_practice

import (
	"context"
	"fmt"
	_"github.com/antlabs/timer"
	"time"
)

// 实现超时退出的方式
/*
	1、context.WithTimeout/context.WithDeadline + time.After
	2、context.WithTimeout/context.WithDeadline + time.NewTimer
	3、channel + time.After/time.NewTimer
 */

// 1、context.WithTimeout/context.WithDeadline + time.After
func TryTimeout1() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second * 2))

	// 有个疑问就是这个
	defer cancel()

	// 模拟调用rpc/http请求操作
	go func(ctx context.Context) {
		// 慢操作
		time.Sleep(time.Second * 3)
		fmt.Println("goroutine done!!")

	}(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("请求成功返回")
		return
	case <-time.After(time.Second):
		fmt.Println("程序超时！退出")
		return
	}

}


// 2、context.WithTimeout/context.WithDeadline + time.NewTimer
func TryTimeout2() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	timer := time.NewTimer(time.Duration(time.Millisecond * 900))
	go func() {
		// http请求
	}()

	select {
	case <-ctx.Done():
		timer.Stop()
		timer.Reset(time.Second)
		fmt.Println("call successfully!!!")
		return
	case <-timer.C:
		fmt.Println("timeout!!!")
		return
	}


}

// 3、channel + time.After/time.NewTimer
func TryTimeout3() {

	done := make(chan struct{}, 1)

	go func() {
		// 模拟处理http请求
		time.Sleep(time.Second)
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(800 * time.Millisecond)):
		fmt.Println("timeout!!!")
		return
	}
}

// NewTimer 创建一个 Timer，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
// t.Reset()需要重置Reset 使 t 重新开始计时
func TryTimeout4() {
	timer := time.NewTimer(time.Second * 2)
	stopC := make(chan struct{})

	go func(t *time.Timer) {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				fmt.Println("timer running")
				t.Reset(time.Second * 2)
			case <-stopC:
				fmt.Println("timer stop!")
			}
		}
	}(timer)

	time.Sleep(time.Second * 10)
	stopC <- struct{}{}
	close(stopC)
	time.Sleep(time.Second * 1)
	fmt.Println("main goroutine退出")
}

// NewTicker 返回一个新的 Ticker，该 Ticker 包含一个通道字段，并会每隔时间段 d 就向该通道发送当时的时间。
//它会调整时间间隔或者丢弃 tick 信息以适应反应慢的接收者。如果d <= 0会触发panic。
//ticker.Stop()关闭该 Ticker 可以释放相关资源。
func TryTimeout5() {
	ticker := time.NewTicker(time.Second * 2)
	stopC := make(chan bool)

	go func(t *time.Ticker) {
		defer t.Stop()

		for {
			select {
			case <-t.C:
				fmt.Println("ticker running!")

			case stop := <-stopC:
				if stop {
					fmt.Println("子goroutine退出，timeout!")
					return
				}
			}
		}


	}(ticker)

	time.Sleep(time.Second * 10)
	stopC <-true
	time.Sleep(time.Second)
	fmt.Println("main goroutine 退出！")

}



func SomeWaysTimeout() {

	//TryTimeout2()

	// 创建一个子节点的context,3秒后自动超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second* 3 )
	defer cancel()
	go watch(ctx, "监控1")
	go watch(ctx, "监控2")


	fmt.Println("现在开始等待8秒,time=", time.Now().Unix())
	time.Sleep(8 * time.Second)

	fmt.Println("等待8秒结束,准备调用cancel()函数，发现两个子协程已经结束了，time=", time.Now().Unix())

}

// 单独的监控协程
func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "收到信号，监控退出,time=", time.Now().Unix())
			return
		default:
			fmt.Println(name, "goroutine监控中,time=", time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}

}