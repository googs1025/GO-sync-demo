package timeout_practice

import (
	"fmt"
	"time"
)

/*
    https://www.cnblogs.com/f-ck-need-u/p/9994512.html
	select如果不去设置超时重试，会容易发生永久阻塞(死锁)的问题。

    1. 可以借助time包的After()实现。
    2. After(d)是只等待一次d的时长，并在这次等待结束后将当前时间发送到通道。Tick(d)则是间隔地多次等待，每次等待d时长，并在每次间隔结束的时候将当前时间发送到通道。
	因为Tick()也是在等待结束的时候发送数据到通道，所以它的返回值是一个channel，从这个channel中可读取每次等待完时的时间点。
 */

func SelectTimeout() {
	//UseTimeAfter1()
	//UseTimeTick1()
	UseTimeTick2()
	//UseTimeTick3()
}


func UseTimeAfter() {
	fmt.Println(time.Now())
	// a为一个只读chan
	a := time.After(time.Second)
	fmt.Println(<-a)
	fmt.Println(a)

}

// 使用After()，也保证了select一定会选中某一个case，这时可以省略default块。
// After()放在select的内部和放在select的外部是完全不一样的，更助于理解的示例见下面的Tick()
func UseTimeAfter1() {
	ch := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep( time.Second)
			ch <-"aaaaaaaa"
		}
		close(ch)

	}()

	for {

		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("消息发送完毕")
				fmt.Println("主goroutine退出")
				return
			}
			fmt.Println("receive the channel value: ", val)
		case <-time.After(time.Millisecond * 500):
			fmt.Println("waiting for the value")
		}

	}

}

func UseTimeTick1() {

	select {
	case <-time.After(7 * time.Second):
		fmt.Println("time After:", time.Now())
	case <-time.Tick(2 * time.Second):
		fmt.Println("time tick:", time.Now())
	}

}


// 特别注意 time.After(4 * time.Second)在select内与外的区别！！

func UseTimeTick2() {


	for {
		// 会一直执行，但是永远不会执行After那行
		select {
		case <-time.After(4 * time.Second):
			fmt.Println("time After:", time.Now().Second())
		case <-time.Tick(2 * time.Second):
			fmt.Println("time tick:", time.Now().Second())
		}
	}


}

func UseTimeTick3() {
	afterC := time.After(4 * time.Second)
	tickC := time.Tick(2 * time.Second)

	for {
		select {
		case <-afterC: // 这里就只会执行一次。
			fmt.Println("time After:", time.Now().Second())
		case <-tickC:	// 会一直执行
			fmt.Println("time tick:", time.Now().Second())
		}
	}

}

