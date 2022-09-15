package channel_mode

import (
	"fmt"
	"golanglearning/new_project/for-gong-zhong-hao/Concurrent-practice/channel-mode/pubsub"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

/*
	发布者订阅者模式：
	消息生产者：publisher
	消息消费者：subscriber
	生产者和消费者是 M:N 的关系
 */

func ChanDemo2()  {

	//
	p := pubsub.NewPublisher(100*time.Second, 10)
	defer p.Close()

	// 订阅全部
	all := p.Subscribe()
	// 订阅字符串上有golang
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})




	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		defer wg.Done()
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		p.Publish("hello, world!")
		p.Publish("hello golang")
	}

	stopC := make(chan os.Signal, 1)


	wg.Wait()

	signal.Notify(stopC, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("收到控制台关闭通知", <-stopC)



}
