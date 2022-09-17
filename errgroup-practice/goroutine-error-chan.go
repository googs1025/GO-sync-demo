package error_practice

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

/*
	实现类似errgroup的功能
	初始化一个chan error类型的chan，然后每个goroutine在遇到错误时候，
	将error写入chan，这样主goroutine通过for range去遍历这个chan就行了。
 */

var wg sync.WaitGroup

func ErrorPractice() {
	num := 5
	TryUseChanAndErr(num)
}

func TryUseChanAndErr(num int) {


	errChan := make(chan error, num)
	wg.Add(num)
	for i :=0 ; i<num; i++ {
		go func(i int) {
			defer wg.Done()
			str := "err" + strconv.Itoa(i)
			errChan <- errors.New(str)
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Println(err)
	}

	fmt.Println("主goroutine退出")



}
