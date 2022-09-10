package pubsub

import (
	"sync"
	"time"
)

// 订阅者是一个chan
type subscriber chan interface{}
// 主题是个过滤器
type topicFunc func(v interface{}) bool

// 发布者对象
type Publisher struct {
	rmu sync.RWMutex
	size int	// chan 大小
	timeout time.Duration	// 超时时间
	subscribers map[subscriber]topicFunc	// 订阅者

}

// 构建函数
func NewPublisher(publishTimeout time.Duration, size int) *Publisher {
	return &Publisher{
		size: size,
		timeout: publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 支持订阅某主题，使用topic过滤
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.size)
	p.rmu.Lock()
	p.subscribers[ch] = topic
	p.rmu.Unlock()

	return ch
}

// 调用Subscribe，就是订阅全部主题
func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 取消订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// 发送主题
func (p *Publisher) Publish(v interface{}) {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)

	}
	wg.Wait()

}

// 发送主题方法。
func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	// 过滤，如果是false，代表不属于这个订阅者关心的内容，直接返回。
	if topic != nil && !topic(v) {
		return
	}
	// 放入chan or 超时
	select {
	case sub <-v:
	case <-time.After(p.timeout):
	}
}

// 关闭订阅通知服务
func (p *Publisher) Close() {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}