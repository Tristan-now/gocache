package geecache

import (
	"fmt"
	"log"
	"sync"
)

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	//首先查找本地缓存有没有key,有则缓存命中
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	//缓存未命中则向其他服务端发起请求
	return g.load(key)
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		//查询group的PeerPicker,用PickPeer方法，根据key去寻找对应的PeerGetter,PeerGetter有Get方法，输入group和key返回value
		//初始化group时，RegisterPeer(peers)为gee添加了HTTPPool这个PeerPicker
		//g.peers.PickPeer(key) 查找 HTTPPool的Map,在哈希环映射里，通过key返回一个PeerGetter（实现了Get方法）
		if peer, ok := g.peers.PickPeer(key); ok {
			// g.getFromPeer从对应的PeerGetter ,调用 peerGetter的Get方法，从gee.name和key拼接出url
			// u := fmt.Sprintf(
			//		"%v%v/%v",
			//		h.baseURL,
			//		url.QueryEscape(group),
			//		url.QueryEscape(key),
			//	)
			// 即是u := peer.baseURL /score/key

			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}
	//从数据库加载key,返回value，加入本地缓存
	return g.getLocally(key)
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
