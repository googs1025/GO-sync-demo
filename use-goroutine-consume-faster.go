package main

import (
	"context"
	"sync"
	"time"
)

func ConsumeTask(ctx context.Context) {
	LOOP:
		var total int
		var success int
		start := time.Now()
		wg := sync.WaitGroup{}
		gLock := sync.Mutex{}
		taskChan := make(chan Task, 50)
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				// 获取的长度是0或者错误 直接break

				// 生产 遇到错误continue

				// 反序列化

				total += 1
				taskChan <- task
			}

			// 结束生产
			close(taskChan)
		}()

		// 多个消费者
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					// 消费数据
					if task, ok := <-taskChan; ok {
						// 。。。。。。
						gLock.Lock()
						success += 1
						gLock.Unlock()
					}
				} else {// chan关闭了 就退出消费
					break
				}
			}
		}()
	}

	log.Warn(ctx, fmt.Sprintf("消费中"))
	wg.Wait()
	log.Warn(ctx, fmt.Sprintf("消费结束"))

	if success == 0 || total == 0 {
		log.Warn(ctx, fmt.Sprintf("当前无待消耗的任务, sleep 10s"))
		time.Sleep(10 * time.Second)
		goto LOOP
	}


	larkText := requestcommon.NewLarkCustomBotContentRichText("消费", time.Now().Format("2006-01-02 15:04:05"))
	totalText := fmt.Sprintf("总共待消费：%d", total)
	failText := fmt.Sprintf("失败：%d", total-success)
	successText := fmt.Sprintf("成功: %d", success)
	takeText := fmt.Sprintf("耗时: %v", time.Since(start))
	ipText := fmt.Sprintf("IP: %s", common.LocalIP())

	// 增加各种指标预警
	larkText.AddTextWithTag(totalText).AddTextWithTag(successText).AddTextWithTag(failText).AddTextWithTag(takeText).AddTextWithTag(ipText)

	// 通过飞书hook url 发送出去
	err := requestcommon.SendLarkCustomBotMsgRichText(ctx, "hook_url", *larkText)
	log.Warn(ctx, fmt.Sprintf("飞书发送消费通知 err: %v", err))

	goto LOOP
}

