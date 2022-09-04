package redis_practice

import "fmt"

// 类似waitgroup
func waitJobs(dones <-chan struct{}, results chan string, jobNum int) {

	working := jobNum
	done := false

	for {
		select {
		case result := <-results:
			fmt.Println(result)
		case <-dones:
			working--
			// 说明已经结束！
			if working <= 0 {
				done = true
			}
		default:
			if done{
				return
			}
		}
	}


}
