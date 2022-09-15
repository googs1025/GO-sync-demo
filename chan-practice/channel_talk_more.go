package chanpractice

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)


// https://mp.weixin.qq.com/s/sIMVKgGgEC0F6ftZjwnn6w



func chantalkmore() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	// 定义好消费者与生产者
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	stopCh := make(chan struct{})
	// stopCh 是额外引入的一个信号 channel.
	// 它的生产者是下面的 toStop channel，
	// 消费者是上面 dataCh 的生产者和消费者
	toStop := make(chan string, 1)
	// toStop 是拿来关闭 stopCh 用的，由 dataCh 的生产者和消费者写入
	// 由下面的匿名中介函数(moderator)消费
	// 要注意，这个一定要是 buffered channel （否则没法用 try-send 来处理了）

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					// try-send 操作
					// 如果 toStop 满了，就会走 default 分支啥也不干，也不会阻塞
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}


				// try-receive 操作，尽快退出
				// 如果没有这一步，下面的 select 操作可能造成 panic
				select {
				case <- stopCh:
					return
				default:
				}

				// 如果尝试从 stopCh 取数据的同时，也尝试向 dataCh
				// 写数据，则会命中 select 的伪随机逻辑，可能会写入数据
				select {
				case <- stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				// 同上
				select {
				case <- stopCh:
					return
				default:
				}

				// 尝试读数据
				select {
				case <- stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
