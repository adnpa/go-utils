package local

import (
	"container/list"
	"fmt"
)

// LRUCache 结构体
type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	order    *list.List
}

// Node 结构体用于存储缓存条目
type Node struct {
	key   int
	value int
}

// NewLRUCache 创建一个新的 LRUCache
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		order:    list.New(),
	}
}

// Get 从缓存中获取值
func (lru *LRUCache) Get(key int) int {
	if elem, found := lru.cache[key]; found {
		// 移动到链表的前面
		lru.order.MoveToFront(elem)
		return elem.Value.(*Node).value
	}
	return -1 // 返回 -1 表示未找到
}

// Put 将键值对放入缓存
func (lru *LRUCache) Put(key int, value int) {
	if elem, found := lru.cache[key]; found {
		// 更新值并移动到链表前面
		elem.Value.(*Node).value = value
		lru.order.MoveToFront(elem)
	} else {
		// 添加新节点
		if lru.order.Len() >= lru.capacity {
			// 移除最少使用的元素
			lru.removeOldest()
		}
		node := &Node{key, value}
		newElem := lru.order.PushFront(node)
		lru.cache[key] = newElem
	}
}

// removeOldest 移除最少使用的元素
func (lru *LRUCache) removeOldest() {
	oldest := lru.order.Back()
	if oldest != nil {
		lru.order.Remove(oldest)
		node := oldest.Value.(*Node)
		delete(lru.cache, node.key)
	}
}

// 示例
func main() {
	cache := NewLRUCache(2) // 创建一个容量为 2 的 LRU 缓存

	cache.Put(1, 1)           // 缓存是 {1=1}
	cache.Put(2, 2)           // 缓存是 {1=1, 2=2}
	fmt.Println(cache.Get(1)) // 返回 1
	cache.Put(3, 3)           // 该操作会使得关键字 2 作废，缓存是 {1=1, 3=3}
	fmt.Println(cache.Get(2)) // 返回 -1 (未找到)
	cache.Put(4, 4)           // 该操作会使得关键字 1 作废，缓存是 {4=4, 3=3}
	fmt.Println(cache.Get(1)) // 返回 -1 (未找到)
	fmt.Println(cache.Get(3)) // 返回 3
	fmt.Println(cache.Get(4)) // 返回 4
}
