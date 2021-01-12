package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

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
	queue    List
	items    map[Key]*ListItem
	mu *sync.Mutex // Для блокировки доступа
}

func (lC *lruCache) Set(key Key, value interface{}) bool {
	lC.mu.Lock()
	defer lC.mu.Unlock()

	if lC.capacity == 0 {
		return false
	}

	// Если элемент есть в словаре
	if _, ok := lC.items[key]; ok {
		lC.items[key].Value.(*cacheItem).value = value

		lC.queue.MoveToFront(lC.items[key])
		lC.items[key] = lC.queue.Front()
		return true
	}

	// Если размер очереди больше емкости кэша, удалаяем последний элемент
	if lC.queue.Len() >= lC.capacity {
		delete(lC.items, lC.queue.Back().Value.(*cacheItem).key)
		lC.queue.Remove(lC.queue.Back())
	}

	// Элемента нет в словаре
	lC.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	lC.items[key] = lC.queue.Front()

	return false
}

func (lC *lruCache) Get(key Key) (interface{}, bool) {
	lC.mu.Lock()
	defer lC.mu.Unlock()

	// Если элемент есть в словаре
	if _, ok := lC.items[key]; ok {
		lC.queue.MoveToFront(lC.items[key])
		lC.items[key] = lC.queue.Front()
		return lC.items[key].Value.(*cacheItem).value, ok
	}

	// Элемента нет в словаре
	return nil, false
}

func (lC *lruCache) Clear() {
	lC.mu.Lock()
	defer lC.mu.Unlock()

	for lC.queue.Back() != nil {
		delete(lC.items, lC.queue.Back().Value.(*cacheItem).key)
		lC.queue.Remove(lC.queue.Back())
	}
}

func NewCache(capacity int) Cache {
	// Если глубина отрицательная, делаем ее равной 0
	if capacity < 0 {
		capacity = 0
	}

	lC := lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    map[Key]*ListItem{},
		mu: &sync.Mutex{},
	}

	return &lC
}
