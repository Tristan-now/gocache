package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// g.mu 这个锁的调用时机：g.m 加入 / 删除 新的key-call映射对时
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	//上锁，判断是否需要加新的映射对
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		// 映射已经存在，解锁
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	// 映射加入完毕，解锁
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	// 删除映射对
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err

}
