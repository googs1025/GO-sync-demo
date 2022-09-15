package pool_practice

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 任务对象
type TaskObj struct {
	// 可以加入对象字段

}

type TaskObjPool struct {
	pool chan *TaskObj
	poolObjNum int
	poolinit bool
}

func NewTaskObjPool(taskNum int, init bool) *TaskObjPool {
	taskPool := &TaskObjPool{
		pool: make(chan *TaskObj, taskNum),
		poolObjNum: taskNum,
		poolinit: init,
	}
    if taskPool.poolinit {

		for i := 0; i < taskPool.poolObjNum; i++ {
			taskPool.pool <- &TaskObj{}
		}

	}


	return taskPool

}

var once sync.Once

func (t *TaskObjPool) GetObj(timeout time.Duration) (*TaskObj, error) {

	poolinit := func() {
		for i := 0; i < t.poolObjNum; i++ {
			t.pool <- &TaskObj{}
		}
	}
	if !t.poolinit {

		once.Do(poolinit)

	}

	select {
	case obj := <-t.pool:
		return obj, nil
	case <-time.After(timeout):
		return nil, errors.New("已经超时退出！")
	}

}

func (t *TaskObjPool) PutObj(obj *TaskObj) error {

	select {
	case t.pool <-obj:
		fmt.Println("成功放回池中！")
		return nil
	default:
		return errors.New("池中已满，返回错误！")
	}

}


func PoolPractice2() {

	pool := NewTaskObjPool(10, false)

	obj1, err := pool.GetObj(time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(obj1)

	err = pool.PutObj(obj1)
	if err != nil {
		fmt.Println(err)
	}


	for i := 0; i < 11; i++ {
		if obj, err := pool.GetObj(time.Second * 1); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(obj)
			if err := pool.PutObj(obj); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	fmt.Println("主goroutine完成退出")

}
