package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {




	signalHandler() //
	//SetupSignalHandler(shutdown)
	select {}  // 用select{} 会让进程hang在这里，只能用"kill -9" 删除！！




}

// 可以制定要监听的信号参数：如 signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM

func signalHandler() {
	// 接受一个signal信号的channel
	signalChan := make(chan os.Signal, 1)
	// 用os库中的Notify 可以接受os传来的信号，并放到chan中
	signal.Notify(signalChan)

	// 采用一个goroutine监听chan
	go func() {
		for {
			select {

			case sig := <-signalChan:
				log.Printf("收到进程给的信号[signal = %v ]", sig)
				return

			default:
				log.Println("尚未收到信号")
				time.Sleep(time.Second)

			}


		}
	}()
}


/*
    **优雅退出进程！**
	退出程序时，通常不能简单粗暴的直接kill，
	需要做一些退出前处理（例如关掉数据库连接，持久化日志，清理应用垃圾等），
    而退出处理的触发就是通过信号接收处理。

 */



// 可以注册不同shutdownFunc，实现不同逻辑
func SetupSignalHandler(shutdownFunc func(bool)) {
	// 接受os通知的chan
	closeSignalChan := make(chan os.Signal, 1)
	// 信号通知函数
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	// 启一个goroutine通知
	go func() {
		sig := <-closeSignalChan
		log.Printf("收到信号[signal = %v ]", sig)
		// 判断退出信号的种类
		shutdownFunc(sig == syscall.SIGQUIT)
	}()


}

// 在程序退出之前，相关资源得到妥善地处理。
func shutdown(isgraceful bool) {

	// 查看通知种类，判断是否优雅退出

	if isgraceful {

		//当满足 sig == syscall.SIGQUIT,做相应退出处理

	}

	// 不是syscall.SIGQUIT的退出信号时，做相应退出处理

}


