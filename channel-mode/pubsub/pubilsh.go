package pubsub

import (
	"sync"
	"time"
)

type subscriber chan interface{}
type topicFunc func(v interface{}) bool


type Publisher struct {
	rmu sync.RWMutex
	size int
	timeout time.Duration
	subscribers map[subscriber]topicFunc
}

func NewPublisher(publishTimeout time.Duration, size int) *Publisher {
	return &Publisher{
		size: size,
		timeout: publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.size)
	p.rmu.Lock()
	p.subscribers[ch] = topic
	p.rmu.Unlock()

	return ch
}

func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) Evict(sub chan interface{}) {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

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

func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <-v:
	case <-time.After(p.timeout):
	}
}

func (p *Publisher) Close() {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}