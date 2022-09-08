package main

import (
	"fmt"
	pubsub "golanglearning/new_project/for-gong-zhong-hao/Concurrent-practice/channel-mode/pubsub"
	"strings"
	"sync"
	"time"
)


func main()  {

	p := pubsub.NewPublisher(100*time.Second, 10)
	defer p.Close()

	all := p.Subscribe()
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


	wg.Wait()



}
