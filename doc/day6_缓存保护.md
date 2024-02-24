> **缓存雪崩**：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。

> **缓存击穿**：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。

> **缓存穿透**：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

通过两个结构体 `call` 和 `Group` 完成缓存保护，使得同一时间对一个key的大量请求合并为一次请求

```
type call struct{
    wg sync.WaitGroup
    val interface{}
    err error
}

type Group struct{
    mu sync.Mutex
    m map[string]*call
}

```

g.mu 锁的作用时机：
- 涉及`g.m[string]*call` 的添加和删除时，要进行加锁