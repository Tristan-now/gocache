回顾分布式缓存获取key的判断流程


- 在当前节点存在对应key缓存：返回value
- 当前节点不存在对应key，使用一致性哈希算法检查是否应从远程节点获取

    - 应从远程节点获取：本地节点的客户端向远程节点的服务端发起请求，获得value
    - 不从远程节点获取：调用回调函数，从数据库从获得value并添加到缓存中
  
```
// Overall flow char										     requsets					        local
// gee := createGroup() --------> /api Service : 9999 ------------------------------> gee.Get(key) ------> g.mainCache.Get(key)
// 						|											^					|
// 						|											|					|remote
// 						v											|					v
// 				cache Service : 800x								|			g.peers.PickPeer(key)
// 						|create hash ring & init peerGetter			|					|
// 						|registry peers write in g.peer				|					|p.httpGetters[p.hashRing(key)]
// 						v											|					|
//			httpPool.Set(otherAddrs...)								|					v
// 		g.peers = gee.RegisterPeers(httpPool)						|			g.getFromPeer(peerGetter, key)
// 						|											|					|
// 						|											|					|
// 						v											|					v
// 		http.ListenAndServe("localhost:800x", httpPool)<------------+--------------peerGetter.Get(key)
// 						|											|
// 						|requsets									|
// 						v											|
// 					p.ServeHttp(w, r)								|
// 						|											|
// 						|url.parse()								|
// 						|--------------------------------------------

```


9999暴漏给用户，用户通过
curl "http://localhost:9999/api?key=Tom" 寻找key
- 8003首先查找本地缓存，不存在
- 查找哈希环，若结果指向自己，则从数据库load并返回
- 若结果不指向自己，通过哈希环查找对应的服务端，访问服务端
- 服务端查找key,缓存中存在则返回，不存在同样查询哈希环，此时指向了自己，则加载数据库放入缓存，并返回结果