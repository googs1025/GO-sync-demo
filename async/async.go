package async

import "fmt"

/*
	1. 使用方式

	async.Execute(func() {
    	someUint64, err = strconv.ParseUint("7770009", 10, 16)
    	if err != nil {
       		panic(err)
    	}
	})

	2. 这个小框架可以保证goroutine在运行时，不会被没有处理的panic使得程序退出。
 */


// 提供异步安全地执行函数的方法，在可能发生panic时进行恢复。
func Execute(fn func()) {

	go func() {
		defer recoverPanic()
		fn()
	}()

}

// 当goroutine发送panic，将错误信息打印到控制台
func recoverPanic() {

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		fmt.Println(err)
	}

}

