package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLru_Set_WhenItemNewAndCacheNotFull_PushesToFront(t *testing.T) {
	// Arrange
	capacity := 4
	lru := New(capacity)
	lru.Set("1", 1)
	lru.Set("2", "2")
	lru.Set("3", struct{}{})

	key := "4"
	val := 3.5

	// Act
	lru.Set(key, val)

	// Assert
	elemFromMap, ok := lru.itemMap[key]
	assert.True(t, ok)
	assert.NotNil(t, elemFromMap)
	assert.Equal(t, elemFromMap.Value.(*item).value, val)

	front := lru.queue.Front()
	assert.NotNil(t, front)
	assert.Equal(t, front.Value.(*item).value, val)

	assert.Equal(t, lru.queue.Len(), capacity)
}

func TestLru_Set_WhenItemExistsAndCacheNotFull_MovesExistingToFront(t *testing.T) {
	// Arrange
	capacity := 4
	lru := New(capacity)
	lru.Set("1", 1)
	lru.Set("2", "2")
	lru.Set("3", struct{}{})

	key := "3"
	val := 3.5

	// Act
	lru.Set(key, val)

	// Assert
	elemFromMap, ok := lru.itemMap[key]
	assert.True(t, ok)
	assert.NotNil(t, elemFromMap)
	assert.Equal(t, elemFromMap.Value.(*item).value, val)

	front := lru.queue.Front()
	assert.NotNil(t, front)
	assert.Equal(t, front.Value.(*item).value, val)

	assert.Equal(t, lru.queue.Len(), 3)
}

func TestLru_Set_WhenItemNewAndCacheFull_RemovesItemAndPushesToFront(t *testing.T) {
	// Arrange
	capacity := 3
	leastRecentKey := "1"
	expectedBack := item{key: "2", value: "2"}
	lru := New(capacity)
	lru.Set(leastRecentKey, 1)
	lru.Set(expectedBack.key, expectedBack.value)
	lru.Set("3", struct{}{})

	key := "4"
	val := 3.5

	// Act
	lru.Set(key, val)

	// Assert
	elemFromMap, ok := lru.itemMap[key]
	assert.True(t, ok)
	assert.NotNil(t, elemFromMap)
	assert.Equal(t, elemFromMap.Value.(*item).value, val)

	_, ok = lru.itemMap[leastRecentKey]
	assert.False(t, ok)

	front := lru.queue.Front()
	assert.NotNil(t, front)
	assert.Equal(t, front.Value.(*item).value, val)

	backElem := lru.queue.Back()
	backItem := backElem.Value.(*item)
	assert.NotNil(t, backElem)
	assert.Equal(t, *backItem, expectedBack)

	assert.Equal(t, lru.queue.Len(), 3)
}

func TestLru_Set_WhenItemExistsAndCacheFull_KeepsAllItemsAndMovesExistingToFront(t *testing.T) {
	// Arrange
	capacity := 3
	lru := New(capacity)
	lru.Set("1", 1)
	lru.Set("2", "two")
	lru.Set("3", struct{}{})

	key := "1"
	val := 3.5

	// Act
	lru.Set(key, val)

	// Assert
	elemFromMap, ok := lru.itemMap[key]
	assert.True(t, ok)
	assert.NotNil(t, elemFromMap)
	assert.Equal(t, elemFromMap.Value.(*item).value, val)

	front := lru.queue.Front()
	assert.NotNil(t, front)
	assert.Equal(t, front.Value.(*item).value, val)

	back := lru.queue.Back()
	assert.NotNil(t, back)
	backItem := back.Value.(*item)
	assert.Equal(t, backItem.key, "2")
	assert.Equal(t, backItem.value, "two")

	assert.Equal(t, lru.queue.Len(), 3)
}
