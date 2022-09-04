package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// https://mp.weixin.qq.com/s/1OoeDckI9Kg2yDLZ_hRRTA

func main() {
	fmt.Println("start")
	//someHanderCancel()
	someHanderTimeout()
	fmt.Println("end")
}


// 例1. 通过取消函数WithCancel()传播
func someHanderCancel() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	//
	for i := 0 ; i < 5; i++ {
		go doSth(ctx, "child goroutine" + strconv.Itoa(i))
	}

	time.Sleep(time.Second * 3)
	cancel()
	time.Sleep(2 * time.Second)

}

func doSth(ctx context.Context, name string) {
	i := 1
	for {

		time.Sleep(time.Second * 5)
		select {
		case <-ctx.Done():
			fmt.Printf("%s done!\n", name)
			fmt.Printf("%s 退出\n", name)
			return
		default:
			fmt.Printf("%s had worked %d seconds \n", name, i)

		}
		i++
	}


}

// 例2. 通过超时控制函数传播
func someHanderTimeout() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3 * time.Second)

	for i := 0 ; i < 5; i++ {
		go doSth(ctx, "child goroutine" + strconv.Itoa(i))
	}

	time.Sleep(time.Second * 10)
	cancel()
	time.Sleep(2 * time.Second)

}


