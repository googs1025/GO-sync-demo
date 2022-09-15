package chanpractice

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 使用buffered channel实现异步处理请求的示例。
/*
	有(最多)3个worker，每个worker是一个goroutine，它们有worker ID。
	每个worker都从一个buffered channel中取出待执行的任务，每个任务是一个struct结构，包含了任务id(JobID)，当前任务的队列号(ID)以及任务的状态(worker是否执行完成该任务)。
	在main goroutine中将每个任务struct发送到buffered channel中，这个buffered channel的容量为10，也就是最多只允许10个任务进行排队。
	worker每次取出任务后，输出任务号，然后执行任务(run)，最后输出任务id已完成。
	每个worker执行任务的方式很简单：随机睡眠0-1秒钟，并将任务标记为完成。
 */


type Task struct {
	ID	int
	JobID int
	Status string
	CreateTime time.Time

}

func (t *Task) run() error {

	// 业务逻辑
	// 或是业务逻辑的方法可以在这调用！
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	t.Status = "Completed"

	return nil
}

func worker(in <-chan *Task, workID int, wg2 *sync.WaitGroup) {

	defer wg2.Done()

	for v := range in {
		fmt.Printf("Worker%d: recv a request: TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
		if err := v.run(); err != nil {
			fmt.Println("some task do somthing wrang!")
			break
		}
		fmt.Printf("Worker%d: Completed for TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
	}

}

const WorkerNum = 3

func chanpractice2() {

	var wg2 sync.WaitGroup
	// 放入的工作对列
	taskqueue := make(chan *Task, 10)

	wg2.Add(3)
	for workID := 0; workID <= WorkerNum; workID++ {
		// 启的工作goroutine
		go worker(taskqueue, workID, &wg2)
	}

	// 生产者的所有数据
	for i := 0; i < 15; i++ {
		taskqueue <-&Task{
			ID: i,
			JobID: 100 + i,
			CreateTime: time.Now(),
		}
	}
	// 生产后记得要关闭！
	close(taskqueue)
	wg2.Wait()


}