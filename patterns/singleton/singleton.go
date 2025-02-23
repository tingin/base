package singleton

import (
	"sync"
)

type SingletonMap[K comparable, V any] struct {
	instances sync.Map
}

func NewSingletonMap[K comparable, V any]() *SingletonMap[K, V] {
	return &SingletonMap[K, V]{}
}

func (s *SingletonMap[K, V]) AddFactory(key K, factory *V) {
	s.instances.Store(key, factory)
}

func (s *SingletonMap[K, V]) Remove(key K) {
	s.instances.Delete(key)
}

func (s *SingletonMap[K, V]) GetInstance(key K) *V {
	if loadedInstance, ok := s.instances.Load(key); ok {
		return loadedInstance.(*V)
	}
	return nil
}
