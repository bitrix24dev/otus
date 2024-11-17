package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.items[key]; exists {
		// Если элемент уже есть в кэше, обновляем его значение и перемещаем в начало
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	// Если элемент новый, добавляем его в начало
	newItem := &cacheItem{key: key, value: value}
	listItem := c.queue.PushFront(newItem)
	c.items[key] = listItem

	// Если кэш переполнен, удаляем последний элемент
	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		if lastItem != nil {
			c.queue.Remove(lastItem)
			delete(c.items, lastItem.Value.(*cacheItem).key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, exists := c.items[key]; exists {
		// Если элемент найден, перемещаем его в начало и возвращаем значение
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Очищаем кэш
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
