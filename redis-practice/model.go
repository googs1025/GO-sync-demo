package redis_practice

import "gopkg.in/redis.v3"

// 任务对象的字段
type Job struct {
	ID string
	Client *redis.Client
	Result chan <-string
}


