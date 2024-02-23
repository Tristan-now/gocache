package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630

$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	//startCacheServer(addrMap[port], addrs, gee)
	//startCacheServer("http://localhost:8001",["http://localhost:8001","http://localhost:8002", "http://localhost:8003"], gee)

	// peers 的 HTTPPool self = addr,即 "http://localhost:8001" ,basePath = "/_geecache/"
	peers := geecache.NewHTTPPool(addr)
	// Set , 为HTTPPool设置 key到httpGetter 的 map映射
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			//Get方法首先在本地查找缓存，不存在则查找哈希环上对应的服务端，不是自己发起请求，是自己就加载数据库
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	//客户端口9999
	apiAddr := "http://localhost:9999"
	//三个服务端口 8001 8002 8003
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		// 输入客户端口9999和group实例
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)
}
