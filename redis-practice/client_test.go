package redis_practice

import (
	"fmt"
	"gopkg.in/redis.v3"
	"runtime"
	"strconv"
	"testing"
	"time"
)

/*
	redis使用命令教程
	https://www.runoob.com/redis/redis-keys.html
    redis进入docker image命令
	docker exec -it af2853715392 redis-cli


	https://www.programminghunter.com/article/4673736664/


 */

// 初始化客户端
func initClient(opt *redis.Options) *redis.Client {

	client := redis.NewClient(opt)
	//if err := client.FlushDb(); err != nil {
	//	panic(err)
	//}

	return client

}

var (
	// 并发数
	jobNum = runtime.NumCPU()
	//每次写入redis的数量
	//除以 jobnum 为了保证改变了任务数, 总量不变, 便于测试
	procnum = 10000 / jobNum
	poolSize = 10
)

func Test_redis(t *testing.T) {

	// redis配置
	opt := &redis.Options{
		Addr: "localhost:6379",
		DialTimeout: time.Second,
		ReadTimeout: time.Second,
		WriteTimeout: time.Second,
		PoolSize: poolSize,
		//Password: "123456",
		DB: 0,
	}

	start := time.Now()
	fmt.Println("start:", start)

	// 退出前执行 统计一下
	defer func() {
		end := time.Now()
		fmt.Println("end:", end)
		fmt.Println("jobs num:", jobNum, "total items:", jobNum*procnum)
		fmt.Println("total seconds:", end.Sub(start).Seconds())
	}()

	// 任务chan
	jobs := make(chan Job, jobNum)
	// 放入结果chan
	results := make(chan string, jobNum*procnum)
	// 通知job结束chan
	dones := make(chan struct{}, jobNum)

	// 初始化client
	client := initClient(opt)
	defer client.Close()
	
	// 任务函数
	jobfunc := func(client *redis.Client, id string) (string, error) {

		// 结束时通知！
		defer func() {
			dones <- struct{}{}
			fmt.Println("job id:", id, "finish")
		}()

		// 执行！
		for idx := 0; idx < procnum; idx++ {
			key := id + "-" + strconv.Itoa(idx)
			// 写入
			value, err := client.Set(key, time.Now().String(), 0).Result()
			if err != nil {
				return "", err
			}
			fmt.Println("key:", key, " | result:", value, " | error:", err)
		}

		return "ok", nil
	}

	// 把job放入chan中
	go func() {
		for index := 0; index < jobNum; index++ {
			jobs <-Job{strconv.Itoa(index), client, results}
		}
		defer close(jobs)
	}()

	// 拿出job 开始并发执行。
	for j := range jobs {
		go func(job Job) {
			res, err := jobfunc(client, job.ID)
			if err != nil {
				fmt.Println(err)
			}
			job.Result <-res
		}(j)
	}

	// 类似waitgroup
	waitJobs(dones, results, jobNum)


}


