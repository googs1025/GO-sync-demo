package chan_practice

import (
	"math/rand"
)

// 建立多个Producer，并用chan传递data，再用一个函数 for-select merge起来。

func Producer1Demo2() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <-rand.Int()
		}
	}()
	return ch
}

func Producer2Demo2() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <-rand.Int()
		}
	}()
	return ch

}

func MergeProducerDemo2() chan int {

	ch := make(chan int, 30)
	go func() {
		for {
			select {
			case v1 := <-Producer1Demo2():
				ch <- v1
			case v2 := <-Producer2Demo2():
				ch <- v2
			}
		}
	}()

	return ch


}

//
//func main() {
//
//	ch := MergeProducerDemo2()
//	for i := 0; i < 200; i++ {
//		fmt.Println(i, <-ch)
//	}
//
//}
