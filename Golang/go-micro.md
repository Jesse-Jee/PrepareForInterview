# micro
利用微服务架构提供一组服务。          

## micro组件
- API: HTTP网关，使用基于路径的解析将HTTP、json请求动态映射到RPC      
- Auth: 使用JWT令牌和基于规则的访问控制进行身份验证和授权
- Broker: 发布订阅消息的异步通信和发布通知。
- Config: 不需要重启就可对服务级别的配置进行动态配置和secrets管理
- Events: 具有有效消息传递，偏移重放和持久存储的消息流。
- Network: 所有内部请求流量的服务间联网，隔离和路由平面
- Proxy: 用于远程访问和任何外部GRPC请求流量的身份感知代理
- Runtime: 服务生命周期和流程管理，从源到运行的自动构建
- Registry: 丰富功能的元数据集中式服务发现和API端点资源管理器
- Store: 具有TTL过期和持久化存储的键值存储，可保持微服务无状态。

# 各组件原理

# 如何实践


# JWT令牌
Json Web Token (跨域认证解决方案)
## 一般的跨域认证方式
- 客户端把用户名密码发送给服务器     
- 服务器验证通过后，在当前对话session中保存相关数据，如用户角色，登录时间等。
- 服务器向用户返回一个session_id,写入用户cookie。
- 用户随后的每一次请求，都会通过cookie，将session_id传回服务器。
- 服务器收到session_id，找到前期保存的数据，由此得知用户身份。

这样存在一个问题，如果是服务集群，就需要session共享。          
所以一种方式是：
- session持久化
- 服务器端不保存session，所有数据保存在客户端。每次请求发回服务器。JWT就是这种模式。

## JWT原理
服务器认证通过后，发送一个json对象给客户端。后面通信时，客户端都要发回这个对象。服务器完全靠这个对象认定用户身份。为了防止用户篡改信息，
服务器在生成对象是会加上签名。         

## JWT数据结构
很长的字符串，中间用点隔成三部分。分别是：               
- header            
    json对象，包括：         
    - alg:签名算法
    - typ:令牌类型 JWT   
- payload负载                 
    json对象，包括：                  
    - iss: 签发人
    - exp: 过期时间
    - sub: 主题
    - aud: 受众
    - nbf: 生效时间
    - iat: 签发时间
    - jti: 编号               
    默认不加密，不能放私密信息。          
- signature签名               
    - 指定一个密钥，只有服务器知道，使用header里指定的算法，产生签名。         

header和payload使用base64URL算法转成字符串。   
base64中有三个字符：+、/、= 在base64URL中，=被忽略，+被替换成-，/被替换成_.          

```header.payload.signature```        

## JWT使用方式
一种是：放到HTTP请求的authorization字段中                      
另一种是在跨域时，JWT放到POST请求中。   

## JWT特点：
- 默认不加密，但可以加密，对生成的token进行一次加密。
- 不加密情况下，不能讲私密信息写到JWT
- 可以用于认证，也可以传递信息。
- 一旦签发，到期之前都始终有效，无法中途废止。
- JWT本身包含认证信息，一旦泄露。所有人都可以获得此令牌权限。为防盗用，JWT有效时间设置短些。重要的要使用时再次验证。
- 建议使用HTTPS传输。          


# go-micro限流方式     
- micro.WrapClient 包装客户端
- micro.WrapServer 保证服务端
    
使用uber limiter插件通过
```go
    QPS := 100
    micro.WrapHandler(limiter.NewHandlerWrapper(QPS)),    
```

# go-micro注册的实现
- registry
    - etcdRegistry
```go
    type etcdRegistry struct {
    	client  *clientv3.Client
    	options registry.Options
    
    	// register and leases are grouped by domain
    	sync.RWMutex
    	register map[string]register
    	leases   map[string]leases
    }
```    



        

# grpc框架

# grpc调用方式

# grpc过程简介



