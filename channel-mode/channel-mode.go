package channel_mode

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

//

func ChannelMode() {
	WaitForResult()
	FanOut()
	WaitForTask1()
	pooling()
	drop()
	cancellation()
	FanOutSem()
}


// 1.等待结果模式
// 使用没缓存的chan实现 两个goroutine的通知机制。
func WaitForResult() {
	ch := make(chan struct{})

	go func(stopChan chan struct{}) {

		// 可以执行子goroutine的业务逻辑
		time.Sleep(time.Second * 2)

		stopChan <-struct{}{}

		fmt.Println("子goroutine收到退出通知！")

	}(ch)

	// 主进程执行业务逻辑！
	fmt.Println("aaaa")


	// 收到通知
	<-ch
	// 收到 可执行退出逻辑！
	time.Sleep(time.Second * 5)
	fmt.Println("父进程收到通知！")
	fmt.Println("--------------------")

}


// 2.扇出/扇入模式
// 包含多个Goroutine向channel发送数据，要保证数据都能接收到。

func FanOut() {
	children := 2000
	ch := make(chan string, children)

	for i := 0; i < children; i++ {
		go func(child int, dataChan chan string) {
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

			dataChan <-"data" + strconv.Itoa(child)
			// 注意：这里不能用"data" + string(child) 会错的！
			fmt.Println("child : sent signal :", child)

		}(i, ch)
	}

	// 主goroutine使用for循环来接收channel里面的数据。sleep模拟执行的任务。
	for children > 0 {
		res := <-ch
		children--
		fmt.Printf("收到res是 %v\n", res)
		fmt.Println("parent : recv'd signal :", children)

	}
	time.Sleep(time.Second)
	fmt.Println("-----------------------")


}

// 3. 等待任务模式
// 无缓冲chan+子goroutine通过channel接收来自主goroutine发送的数据，也可以是执行任务的函数。

func WaitForTask() {
	// 使用不带缓存的channel，子goroutine等待channel发送数据，接收并执行任务。
	ch := make(chan string)

	// 启一个goroutine 接收数据
	go func(dataChan chan string) {
		// 当chan 没有收到数据时，会阻塞在这
		d := <-dataChan
		fmt.Println("child : recv'd signal :", d)
	}(ch)

	// 业务逻辑
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	// 主goroutine发送数据时
	ch <- "data"
	fmt.Println("parent : sent signal")

	// 执行业务逻辑
	time.Sleep(time.Second)
	fmt.Println("---------------")

}

// 3. 等待任务模式
// 有缓冲chan+子goroutine通过channel接收来自主goroutine发送的数据，也可以是执行任务的函数。

func WaitForTask1() {
	// 有缓冲chan
	ch := make(chan string, 50)

	// 启一个goroutine 接收数据
	go func(dataChan chan string) {
		// 不断接收！
		for {
			select {
			case d := <-dataChan:
				fmt.Println("child : recv'd signal :", d)
				time.Sleep(time.Second * 2)
			default:
				fmt.Println("目前没有收到任务")
				time.Sleep(time.Second)

			}

		}

	}(ch)

	// 可执行业务逻辑
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	// 发送数据！ 可以换成一个函数或方法，作为生产者
	for i :=0 ; i < 1000; i++ {
		ch <- "data" + strconv.Itoa(i)
		time.Sleep(time.Second)

	}
	close(ch) // 发送完记得要close!
	fmt.Println("parent : sent signal")

	time.Sleep(time.Second*50)
	fmt.Println("---------------")

}

// 4. Goroutine池
// 该模式使用了等待任务模式，允许根据资源情况限制子goroutine的个数。
// 主goroutine

func pooling() {
	ch := make(chan string,10)
	// 该函数可以读取机器cpu核数，也就是能并行执行代码的cpu核数
	g := runtime.GOMAXPROCS(0)
	for i := 0; i < g; i++ {
		go func(child int, dataChan chan string) {
			// 这种常用
			for d := range dataChan {
				fmt.Printf("child %d : recv'd signal : %s\n", child, d)
			}

			// 这种读取channel的方式也可以
			//for {
			//	d, exis := <-dataChan
			//	if !exis {
			//		break
			//	}
			//	fmt.Printf("child %d : recv'd signal : %s\n", child, d)
			//}

			fmt.Printf("child %d : recv'd shutdown signal\n", child)

		}(i, ch)
	}

	const work = 100

	for i := 0; i < work; i++ {
		ch <-"data" + strconv.Itoa(i)
		fmt.Println("parent : sent signal :", i)
	}

	close(ch) // 一定要关闭close 可以比较看看区别！
	// 关闭close可以管理之前启的goroutine的生命周期！


	fmt.Println("parent : sent shutdown signal")
	time.Sleep(time.Second)
	fmt.Println("------------------")
}

// 5.Drop模式
// 该模式在写入channel的数据量比较大的时候，超出缓冲的容量就选择丢弃数据。例如当应用程序负载太大就可以丢弃一些请求。

func drop() {
	const cap = 100
	ch := make(chan string, cap)

	go func(dataChan chan string) {

		for i := range dataChan {
			fmt.Println("child : recv'd signal :", i)
		}

	}(ch)

	const work = 2000

	for i := 0; i < work; i++ {
		// 这个select使用了default关键字，它将select转换为非阻塞调用。关键就在这里，如果channel缓冲区满了，select就会执行default。
		select {
		case ch <- "data" + strconv.Itoa(i):
			fmt.Println("parent : sent signal :", i)
		default:
			fmt.Println("parent : dropped data :", i)
		}
	}
	close(ch)
	fmt.Println("parent : sent shutdown signal")
	time.Sleep(time.Second*50)
	fmt.Println("-----------------")

}

// 6.取消模式
// 取消模式用于在执行一些IO操作的时候，可以选择超时时间。你可以选择取消操作，或者直接退出。

func cancellation() {
	ctx := context.Background()
	duration := time.Second * 5
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	ch := make(chan string, 1)

	go func(dataChan chan string) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
		dataChan <-"data"

	}(ch)

	select {
	case d := <-ch:
		fmt.Println("work complete", d)
	case <-ctx.Done():
		fmt.Println("work cancelled")
	}

}

// 可以随时控制可执行的Goroutine数量。
/*
	一开始创建了一个缓冲为2000的channel。和前面的扇入/扇出没啥区别。
	另一个chennel sem也被创建了，在每个子goroutine内部使用，
	可以控制子Goroutine是否能够写入数据容量，缓冲区满的话子goroutine就会阻塞。
	后面的for循环用于等待每个goroutine执行完成。
 */
func FanOutSem() {
	children := 2000
	ch := make(chan string, children)
	g := runtime.GOMAXPROCS(0)
	sem := make(chan bool, g)
	
	for i := 0; i < children; i++ {

		go func(child int) {
			sem <- true

			{
				t := time.Duration(rand.Intn(200)) * time.Millisecond
				time.Sleep(t)
				ch <- "data"
				fmt.Println("child : sent signal :", child)
			}
			<-sem


		}(i)
	}

	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}

	time.Sleep(time.Second)
	fmt.Println("-------------")

}

// 重试超时模式
// 在网络服务中很实用，例如在连接数据库的时候，发起ping操作可能会失败，但是并不希望马上退出，而是在一定时间内发起重试。
func RetryTimeout(ctx context.Context, retryInterval time.Duration, check func(ctx context.Context) error) {
	for {

		fmt.Println("perform user check call")
		if err := check(ctx); err == nil {

			fmt.Println("work finished successfully")
			return

		}
		fmt.Println("check if timeout has expired")
		if ctx.Err() != nil {

			fmt.Println("time expired 1 :", ctx.Err())
			return

		}
		fmt.Printf("wait %s before trying again\n", retryInterval)
		// 创建一个计时器
		t := time.NewTimer(retryInterval)
		select {
		case <-ctx.Done():

			fmt.Println("timed expired 2 :", ctx.Err())
			t.Stop()
			return
		// 定时执行！
		case <-t.C:
			fmt.Println("retry again")
		}
	}
}

// 这个写法有点问题。
// 可以创建一个单独的channel用来实现取消的功能。
func channelCancellation(stop <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-stop:
			cancel()
		case <-ctx.Done():
		}
	}()


	go func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://www.ardanlabs.com/blog/index.xml",
			nil,

		)
		if err != nil {
			return err
		}



		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err

		}
		return nil

	}(ctx)
}
