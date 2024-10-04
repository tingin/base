package singleton

import (
	"fmt"
	"sync"
)

type SingletonMap[K comparable, V any] struct {
	instances sync.Map
	factories sync.Map
	onceMap   sync.Map
}

func NewSingletonMap[K comparable, V any]() *SingletonMap[K, V] {
	return &SingletonMap[K, V]{}
}

func (s *SingletonMap[K, V]) AddFactory(key K, factory func() *V) {
	s.factories.Store(key, factory)
}

func (s *SingletonMap[K, V]) Remove(key K) {
	s.factories.Delete(key)
	s.instances.Delete(key)
	s.onceMap.Delete(key)
}

func (s *SingletonMap[K, V]) GetInstance(key K) *V {
	once, _ := s.onceMap.LoadOrStore(key, &sync.Once{})
	once.(*sync.Once).Do(func() {
		if factory, ok := s.factories.Load(key); ok {
			instance := factory.(func() *V)()
			s.instances.Store(key, instance)
		} else {
			fmt.Printf("Factory for key %v not found!\n", key)
		}
	})
	instance, _ := s.instances.Load(key)
	return instance.(*V)
}
