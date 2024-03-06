分布式缓存系统   （Golang）                                                 	
1.实现单机缓存和基于HTTP的分布式缓存，通过HTTP协议在多个服务器之间共享缓存数据，提高系统的可扩展性和容错能力
2. 基于LRU缓存策略，跟踪数据访问频率，自动移除最长时间未被访问数据，最大化缓存效率和响应时间
3. 使用Go锁防止缓存击穿，使用互斥锁（sync.Mutex）来同步对缓存的访问，防止同时间对相同key的大量重复访问，确保系统稳定性和数据库安全
4. 使用一致性哈希选择节点，实现负载均衡，将节点映射到一个哈希环上，在查询key时只需要去访问对应节点，有效实现负载均衡，提高分布式缓存系统的稳定性和效率
5. 使用protobuf优化节点通信，允许节点之间以紧凑的二进制格式交换数据，减少了网络传输的开销，提高了通信速度和效率，保证了数据的一致性和兼容性
