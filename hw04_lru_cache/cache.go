package hw04_lru_cache //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	List
	items       map[Key]*listItem
	*sync.Mutex // Для блокировки доступа
}

func (lC *lruCache) Set(key Key, value interface{}) bool {
	lC.Lock()
	defer lC.Unlock()

	// Если элемент есть в словаре
	if _, ok := lC.items[key]; ok {
		lC.items[key].Value.(*cacheItem).value = value

		lC.MoveToFront(lC.items[key])
		lC.items[key] = lC.Front()
		return true
	}

	// Если размер очереди больше емкости кэша, удалаяем последний элемент
	if lC.Len() >= lC.capacity {
		delete(lC.items, lC.Back().Value.(*cacheItem).key)
		lC.Remove(lC.Back())
	}

	// Элемента нет в словаре
	lC.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	lC.items[key] = lC.Front()

	return false
}

func (lC *lruCache) Get(key Key) (interface{}, bool) {
	lC.Lock()
	defer lC.Unlock()

	// Если элемент есть в словаре
	if _, ok := lC.items[key]; ok {
		lC.MoveToFront(lC.items[key])
		lC.items[key] = lC.Front()
		return lC.items[key].Value.(*cacheItem).value, ok
	}

	// Элемента нет в словаре
	return nil, false
}

func (lC *lruCache) Clear() {
	lC.Lock()
	defer lC.Unlock()

	lC.List = NewList()
	lC.items = make(map[Key]*listItem, lC.capacity)
}

func NewCache(capacity int) (Cache, error) {
	if capacity < 1 {
		return nil, errors.New("capacity < 1")
	}

	lC := lruCache{
		capacity: capacity,
		List:     NewList(),
		items:    make(map[Key]*listItem, capacity),
		Mutex:    &sync.Mutex{},
	}

	return &lC, nil
}
