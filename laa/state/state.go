package state

import (
	"log"
	"reflect"
	"sync"
	"time"
)

type StateFlow[T any] struct {
	value       T
	mu          sync.RWMutex
	subscribers map[int]chan T
	nextSubID   int
	once        sync.Once
}

func NewGoStateFlow[T any](initialValue T) *StateFlow[T] {
	s := &StateFlow[T]{
		value:       initialValue,
		subscribers: make(map[int]chan T),
		nextSubID:   0,
	}
	return s
}
func (g *StateFlow[T]) Get() T {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}
func (g *StateFlow[T]) Set(newValue T) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if reflect.DeepEqual(g.value, newValue) {
		return
	}
	g.value = newValue
	for _, subChan := range g.subscribers {
		go func(ch chan T, val T) {
			select {
			case ch <- val:
			case <-time.After(1 * time.Second):
				log.Printf("发送给订阅者超时，可能通道已满或订阅者停止监听。值: %v\n\n", val)
			}
		}(subChan, newValue)
	}
}
func (g *StateFlow[T]) Subscribe() (<-chan T, func()) {
	g.mu.Lock()
	defer g.mu.Unlock()
	subID := g.nextSubID
	g.nextSubID++
	subscriberChan := make(chan T, 1)
	select {
	case subscriberChan <- g.value:
	default:
		log.Println("警告: 无法立即发送初始值给新订阅者。")
	}
	g.subscribers[subID] = subscriberChan
	cancelFunc := func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if ch, ok := g.subscribers[subID]; ok {
			delete(g.subscribers, subID)
			close(ch)
			//fmt.Printf("订阅者 ID %d 已取消订阅并关闭通道。\n", subID)
		}
	}
	return subscriberChan, cancelFunc
}
