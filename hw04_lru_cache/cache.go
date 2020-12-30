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
	sync.Mutex	// Для блокировки доступа
	queue List
	items map[Key]*ListItem
}

func (lC *lruCache) Set(key Key, value interface{}) bool {
	if lC.capacity == 0 {
		return false
	}

	lC.Lock()
	defer lC.Unlock()

	// Если элемент есть в словаре
	if _, ok := lC.items[key]; ok {
		lC.items[key].Value.(*cacheItem).value = value

		lC.queue.MoveToFront(lC.items[key])
		lC.items[key] = lC.queue.Front()
		return true
	}

	// Элемента нет в словаре
	lC.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	lC.items[key] = lC.queue.Front()

	// Если размер очереди больше емкости кэша, удалаяем последний элемент
	if lC.queue.Len() > lC.capacity {
		lC.queue.Remove(lC.queue.Back())
	}

	return false
}

func (lC *lruCache) Get(key Key) (interface{}, bool) {
	// Если элемент есть в словаре
	if v, ok := lC.items[key]; ok {
		lC.Lock()
		defer lC.Unlock()

		lC.queue.MoveToFront(v)
		v = lC.queue.Front()
		return v.Value.(*cacheItem).value, ok
	}

	// Элемента нет в словаре
	return nil, false
}

func (lC *lruCache) Clear() {
	lC.Lock()
	defer lC.Unlock()

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

	lC := lruCache {
		capacity: capacity,
		queue:    NewList(),
		items:    map[Key]*ListItem{},
	}

	return &lC
}
