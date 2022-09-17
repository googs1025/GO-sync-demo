package best

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Go并发最佳实践
// https://mp.weixin.qq.com/s/66vI_Zq9Oeb_2IJOga7_qg

func Test() {
	ServerUse()
	NetUse()

}



// 1、使用goroutines管理服务状态
// 通常启动一个http服务器或其他需要长期在后台运行的任务时，
// 可以使用chan或带chan字段的结构体，实现goroutines的同步：

func ServerUse() {
	s := NewServer()
	time.Sleep(time.Second * 10)
	s.Stop()

}

type Server struct {
	// ...其他对象。。。

	// 用来通知结束Chan
	stop chan struct{}
}


// 建立对象
func NewServer() *Server {
	s := &Server{
		stop: make(chan struct{}),
	}
	// 启一个goroutine，用于不断循环检查是否需要结束
	go s.run()
	return s

}

func (s *Server) run() {
	for {
		select {
		case <-s.stop:
			fmt.Println("finishing task")
			time.Sleep(time.Second)
			fmt.Println("task done")

			return

		case <-time.After(time.Second):
			fmt.Println("running task")
		}

	}

}

func (s *Server) Stop() {
	fmt.Println("server stopping")
	s.stop <- struct{}{}
	fmt.Println("server stopped")
}


// 2. 使用buffer channel避免goroutine泄漏


func SendMsg(msg, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = fmt.Fprintf(conn, msg)
	return err
}

// 存在几个问题：
// 写channel会被阻塞
// goroutine保存对chan的引用
// chan无法被垃圾回收

// broadcastSendMsgBad向多个服务发送同一条消息：
func broadcastSendMsgBad(msg string, addrs []string) chan error {
	var wg sync.WaitGroup
	errChan := make(chan error)  // 无缓冲chan
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			errChan <- SendMsg(msg, addr)
			fmt.Println("done 一个任务")
		}(addr)
	}

	wg.Wait()

	return errChan

}

func NetUse() {
	addrs := []string{
		"localhost:8080",
		"http://baidu.com",
		"http://google.com",

	}

	errChan := broadcastSendMsgGood("hello world", addrs)
	//errChan := broadcastSendMsgBad("hello world", addrs) // 用不了，会直接阻塞。
	//errChan := broadcastSendMsg("hello world", addrs)

	// 主goroutine消费者。
	for err := range errChan {
		fmt.Println(err)
	}

	fmt.Println("消息发送成功")
}

// 带缓冲带的通道可以保证每个goroutine都能向通道里面写入数据并结束goroutine
func broadcastSendMsg(msg string, addrs []string) chan error {

	var wg sync.WaitGroup	// 用来waitgroup并发
	// 用来存储err
	errChan := make(chan error, len(addrs))
	defer close(errChan)	// 一个消费者， 需要

	// 启多个goroutine作为生产者。
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			errChan <- SendMsg(msg, addr)
			fmt.Println("done 一个任务")

		}(addr)
	}

	// 会阻塞在这里
	wg.Wait()


	return errChan

}

// 使用stopChan通道来保证所有的goroutine能正常退出，
// 这样就不会存在goroutine泄漏。不管发送消息是否成功还是失败，
// 当broadcastSendMsgGood函数结束时会调用close(stopChan)，保证所有goroutine的退出。
func broadcastSendMsgGood(msg string, addrs []string) chan error {

	var wg sync.WaitGroup
	errChan := make(chan error, len(addrs))
	defer close(errChan)	// 用完要关闭！
	stopChan := make(chan struct{})


	// 启多个goroutine 作为生产者
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			// 不断循环 监听不同chan
			for {
				select {
				case errChan <- SendMsg(msg, addr):
					fmt.Println("done 一个任务")
				case <- stopChan:
					fmt.Println("task done")
					return
				}
			}

		}(addr)
	}


	// 定时触发关闭chan的任务
	time.AfterFunc(time.Second * 5, func() {
		close(stopChan)
	})

	// 阻塞在这里
	wg.Wait()


	return errChan

}

