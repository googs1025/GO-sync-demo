package cron_practice

import (
	"fmt"
	"strconv"
	"time"
)

func TimeTicker() {
	input := make(chan interface{})

	go func(ch chan interface{}) {
		for i := 0; i < 10; i++ {
			value := strconv.Itoa(i)
			time.Sleep(time.Second * 3)
			ch <- "hellp world" + value
		}

		close(ch)
	}(input)

	timer1 := time.NewTicker(time.Second * 10)
	timer2 := time.NewTicker(time.Second * 5)

	for {
		select {
		case msg, ok := <-input:
			if !ok {
				fmt.Println("生产者发送完毕！")
				return
			}
			fmt.Println(msg)
		case <-timer1.C:
			fmt.Println("计时器1 时间到")
			timer1.Reset(time.Second * 10)
		case <-timer2.C:
			fmt.Println("计时器2 时间到")
			timer2.Reset(time.Second * 5)
		}
	}


}
