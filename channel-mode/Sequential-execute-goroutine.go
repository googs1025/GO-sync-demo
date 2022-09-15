package channel_mode

import (
	"fmt"
	"sync"
	"time"
)

/*
    https://www.cnblogs.com/f-ck-need-u/p/9994652.html
	上面的示例中：run1 goroutine被ch1阻塞，run2 goroutine被ch2阻塞，run3 goroutine被ch3阻塞。run3依赖的ch3由B关闭，run2依赖的ch2由run1关闭。

	如此一来，当main goroutine中的x被关闭后，run1()从阻塞中释放，继续执行，关闭ch2，然后run2从阻塞中释放，继续执行，关闭ch3，run3得以释放。由于ch3被关闭后，ch3仍然可读，所以多次执行run3不会出问题。

	run1()和run2()不能多次执行，因为close()不能操作已被关闭的channel。

	注意，上面的channel都是struct{}类型的，整个过程中3个通道都没有传递数据，而是直接关闭来释放通道，让某些阻塞的goroutine继续执行下去。显然，这里的x、y、z的作用都是"信号通道"，用来传递消息。
 */

func Sequential_execute_goroutine() {

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(5)
	go run3(ch3, &wg)
	go run3(ch3, &wg)
	go run2(ch2, ch3, &wg)
	go run1(ch1, ch2, &wg)
	go run3(ch3, &wg)

	close(ch1)
	wg.Wait()


}

func run1(ch1, ch2 chan struct{}, wg *sync.WaitGroup) {

	<-ch1
	fmt.Println("执行run1!!")
	time.Sleep(time.Second)
	close(ch2)
	wg.Done()

}

func run2(ch1, ch2 chan struct{}, wg *sync.WaitGroup) {

	<-ch1
	fmt.Println("执行run2!!")
	time.Sleep(time.Second)
	close(ch2)
	wg.Done()

}

func run3(ch1 chan struct{}, wg *sync.WaitGroup) {

	<-ch1
	fmt.Println("执行run3!!")
	wg.Done()

}
