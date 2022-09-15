package chan_practice

import (
	"fmt"
	"testing"
)

func TestDemo1(t *testing.T) {
	ch := Producer()
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}

func TestDemo2(t *testing.T) {
	ch := MergeProducerDemo2()
	for i := 0; i < 200; i++ {
		fmt.Println(<-ch)
	}
}

func TestDemo3(t *testing.T) {
	doneC := make(chan struct{})
	ch := MergeProducer(doneC)
	for i := 0; i < 200; i++ {
		fmt.Println(<-ch)
	}
	doneC <- struct{}{}
	fmt.Println("主goroutine退出！")
}
