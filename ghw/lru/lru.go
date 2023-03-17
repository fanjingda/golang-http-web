package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64                    //运行使用最大内存
	nbytes   int64                    //当前使用的内存
	ll       *list.List               //队列
	cache    map[string]*list.Element //映射表
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的key
// 好处：淘汰队尾节点时，需要用key从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 用于返回缓存值所占用的内存大小
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	//从字典中查找对应的缓存节点是否存在
	if ele, ok := c.cache[key]; ok {
		//如果存在则将对应节点移动到队首
		c.ll.MoveToFront(ele)
		//返回查找到的值，将entry转为value值类型
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	//取队尾元素
	ele := c.ll.Back()
	//如果存在，则跟新链表和映射表中的对应项
	if ele != nil {
		//把该结点从链表中删除
		c.ll.Remove(ele)
		//从字典汇总删除该节点的映射关系
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//更新当前占用的内存（减去key+value）
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//如果回调函数不为空，则删除元素时调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	//从制度中查找对应的缓存节点是否存在
	if ele, ok := c.cache[key]; ok {
		//存在（跟新）
		//将节点移动到队首
		c.ll.MoveToFront(ele)
		//计算当前占用内存（可以未跟新，只有value跟新）
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		//跟新value值
		kv.value = value
	} else {
		//不存在
		//构建新节点插入到队首
		ele := c.ll.PushFront(&entry{key, value})
		//字典中建立Key-value映射关系
		c.cache[key] = ele
		//计算当前内存（key+value）
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
