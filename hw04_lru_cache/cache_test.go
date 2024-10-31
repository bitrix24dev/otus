package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	/*
		Дополнительный тест на логику выталкивания элементов из-за размера очереди
		 (например: n = 3, добавили 4 элемента - 1й из кэша вытолкнулся);
	*/

	t.Run("push logic when capacity exeed", func(t *testing.T) {
		// Создаем кэш с емкостью 3
		cache := NewCache(3)

		// Добавляем 3 элемента
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		// Проверяем, что все 3 элемента присутствуют в кэше
		val, ok := cache.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", val)

		val, ok = cache.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", val)

		val, ok = cache.Get("key3")
		assert.True(t, ok)
		assert.Equal(t, "value3", val)

		// Добавляем 4-й элемент, что должно вытолкнуть первый элемент ("key1")
		cache.Set("key4", "value4")

		// Проверяем, что "key1" был вытолкнут
		_, ok = cache.Get("key1")
		assert.False(t, ok)

		// Проверяем, что остальные элементы присутствуют
		val, ok = cache.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", val)

		val, ok = cache.Get("key3")
		assert.True(t, ok)
		assert.Equal(t, "value3", val)

		val, ok = cache.Get("key4")
		assert.True(t, ok)
		assert.Equal(t, "value4", val)
	})

	/*
		Дополнительный тест на логику выталкивания давно используемых элементов
		(например: n = 3, добавили 3 элемента, обратились несколько раз к разным элементам:
		изменили значение, получили значение и пр. - добавили 4й элемент,
		из первой тройки вытолкнется тот элемент, что был затронут наиболее давно)
	*/
	t.Run("last recently used", func(t *testing.T) {
		// Создаем кэш с емкостью 3
		cache := NewCache(3)

		// Добавляем 3 элемента
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		// Обращаемся к элементам, чтобы изменить их порядок использования
		cache.Get("key1")                   // Теперь "key1" самый недавно использованный
		cache.Set("key2", "value2_updated") // Обновляем "key2", теперь он самый недавно использованный

		// Добавляем 4-й элемент, что должно вытолкнуть "key3", так как он не использовался
		cache.Set("key4", "value4")

		// Проверяем, что "key3" был вытолкнут
		_, ok := cache.Get("key3")
		assert.False(t, ok)

		// Проверяем, что остальные элементы присутствуют
		val, ok := cache.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", val)

		val, ok = cache.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2_updated", val)

		val, ok = cache.Get("key4")
		assert.True(t, ok)
		assert.Equal(t, "value4", val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
