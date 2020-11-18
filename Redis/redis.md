# Redis支持的数据类型及底层原理
- string
    - SDS 简单动态字符串
- list
    - 双向链表
    - zipList压缩列表
- hash
    - 哈希表
    - 压缩列表
- set
    - 哈希表
    - 整数数组
- sorted set（zset）
    - 压缩列表
    - 跳表
    
拓展的话，还有：
- bitmap 用于二值统计，非0即1的场景，如签到
- Hyperloglog 用于统计基数
- GEO 基于位置信息服务的应用，附近的餐馆等。

## string
### 使用
set key value  
get key  
mset key1 value1 key2 value2  
mget key1 key2  
INCR key 递增数字  
INCRBY key increment 递增指定数字  
DECR key 递减数字   
DECRBY key increment 递减指定数字  
strlen key 获取字符串长度

### 分布式锁
setnx key value   
set key value [EX seconds][PX milliseconds][NX|XX]  
**EX**:key 存在多少秒后过期     
**PX**：key在多少毫秒后过期
**NX**：当key不存在是才能创建key 效果等同于setnx
**XX**：当key存在是覆盖key

### 应用场景
- 商品编号、订单号采用INCR生成
- 点赞数之类的

## hash
### 使用
Hset key filed value  
hget key filed  
hmset key filed1 value1 filed2 value2  
hmget key filed1[filed...]  
hgetall key  
hlen  
hdel  

### 应用场景
购物车

## list
### 使用
LPUSH key value
RPUSH key value
LRANGE key start stop
LLEN key

### 应用场景
微信文章订阅公众号：公众号发布的新文章，将其ID放到我这个关注者的list中来。

## set
### 使用
SADD key member  
SREM key member  删除
SMEMBERS key  获取集合中所有元素  
SISMEMBERS key member  判断元素是否在集合中     
SCARD key  集合中元素个数  
SRANDMEMBER key [数字] 从集合中随机弹出一个元素，不删除    
SPOP key 从集合中随机弹出一个元素，删除。  

### 集合运算
- 差集 SDIFF key
- 交集 SINTER key
- 并集 SUNION key  


### 应用场景
- 微信抽奖小程序 SRANDMEMBER 如果不重复获奖，用SPOP
- 朋友圈点赞 
- 微博好友关注，共同关注的人，我关注的人也关注了他
- 推可能认识的人

## zset
### 使用
ZADD key score member  
ZRANGE key start stop 从小到大排序  
ZSCORE key member  获取元素分数  
ZREM key member  删除元素  
ZRANGEBYSCORE key min max 获取指定范围元素    
ZINCRBY key increment member 增加某个元素分数  
ZCARD key  获取集合中元素数量  
ZCOUNT key min max指定分数范围内的元素个数  

### 应用场景
- 热销
- 打赏排行  



# Redis内存淘汰策略
执行内存淘汰策略前，要先经过删除策略。

## Redis过期键删除策略
- 定时删除：立即删除能保证内存新鲜度，因为它能在键一过期立马删除，占用的内存也会随之释放。
    - 但是，对CPU不友好，因为删除会占CPU时间。实时删除，会让CPU性能损耗，影响数据的读取操作。
- 惰性删除：数据到达过期时间，不处理。等下次访问该数据时，如果没过期，返回数据，如果过期了，删除。
    - 但是，对内存不友好。万一这个数据后面不被访问了呢？一直占用着内存。
- 定期删除：每隔一段时间执行一次删除过期键的策略。并通过限制删除操作执行的时长和频率来减少删除操作对CPU的影响。
    - 周期性轮询时效性数据，采用随机抽取策略，利用过期数据占比的方式控制删除额度。
        - 检测频度可自定义设置
        - 长期占用内存的冷数据会被清除。
    - 难点：确定删除操作执行的时长和频率。太频繁的话CPU一样要崩溃。执行不频繁的话，和惰性一样，还是会导致数据不能及时删掉。
    
**上述三种都有漏洞，必须有兜底的策略。**     


## 8种缓存淘汰策略
在redis.conf中设置的。
默认为不驱逐

- 不驱逐：默认，存满返回错误。我存满了，不干了。
- allkeys-lru：对所有keys使用lru算法删除
- volatile-lru：对所有设置了过期时间的key使用LRU算法删除。
- allkeys-random：对所有key进行随机删除
- volatile-random：对所有设置了过期时间的key随机删除。
- volatile-TTL：删除马上过期的key
- allkeys-lfu：对所有key使用lfu算法删除。
- volatile-lfu：对所有设置过期时间的key使用lfu算法删除。

LRU：最近最少使用
LFU：最近最少频率使用

## 淘汰策略总结
- 两个维度
    - 过期的中选
    - 所有的里面选
- 四个方面
    - LRU
    - LFU
    - random
    - TTL
    
## 你平时用哪一种
allkeys-lru

### 如何配置，如何修改？
- redis.conf文件中设置 maxmemory-policy allkeys-lru
- config set maxmemory-policy allkeys-lru




#zset时间复杂度



#缓存穿透，缓存击穿，雪崩。


#数据结构底层原理及使用



#更新redis缓存与数据库数据不一致问题


#持久化比较
##RDB
##AOF
###AOF重写机制


#hash扩容

#rehash


#Redis线程模型

#如何提高缓存命中