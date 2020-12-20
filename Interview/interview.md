- socket相关知识

- 如何根据struct中的参数对一个struct切片进行排序？                   
实现sort的 Len(),Swap(),Less() 三个方法即可                      

    例如：
```go
    type User struct {
    	Id int
    }
    
    type ById []User
    
    func (b ById) Len() int {
    	return len(b)
    }
    
    func (b ById) Swap(i, j int) {
    	b[i], b[j] = b[j], b[i]
    }
    
    func (b ById) Less(i, j int) bool {
    	return b[i].Id < b[j].Id
    }
    
    func main() {
    	var u1 ById
    	u1 = []User{{1}, {6}, {3}, {5}, {8}, {4}}    
    	sort.SliceStable(u1, u1.Less)
    	fmt.Println(u1)
    }

```


- 负载均衡算法有哪些？
    - 轮询
    - 随机
    - 随机轮询
    - 加权轮询
    - 最少使用
    - 加权最少使用
    - 加权随机
    - 源地址散列
    - 源地址端口散列
    
    
- 常见的MySQL优化方式
    - 硬件优化？
    - 软件升级？ 
    - 减少关联子查询改为关联查询？
    - 适当提升冗余，空间换时间。                     
    
  
    
- docker的隔离如何实现的？               
    使用linux namespace                   
    
- sync.Map的使用
  多线程从map中获取某个key值，返回bool。                      
  var m sync.Map
  if _, ok := m.Load(key);ok {
        return key
  }
  
  如何保存一个key                                         
  m.Store(key, value)                   
  
  如何查询一个value？                  
  m.LoadOrStore(key, value)                         
  如果key存在，就返回value，如果不存在，就保存并返回value。                           
  
   
- gorm的使用
    更新一堆字段？                     
    updates方法                       
    如何执行事务？                     
    db.Transaction()                                    
    支持手动事务                              
    tx := db.Begin()                
    tx.Create()...                              
    tx.Commit()...              
    
    
- 实现一个server，8080端口启动，返回指定json。             
```go
func main() {
	r := gin.Default()
	r.GET("/user", func(context *gin.Context) {
		context.JSON(200, "hello world")
	})
	r.Run()
}
```    

- 几个json,如何设计结构体？                           
{name,pwd}
{name,pwd,info}                         
{name,pwd,info,moreInfo}                

```go
type User struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type UserInfo struct {
	User
	Info string `json:"info"`
}

type UserMore struct {
	UserInfo
	More string `json:"more"`
}

func main() {
	var u = User{Name: "1", Pwd: "1"}
	strU, _ := json.Marshal(u)
	fmt.Println(string(strU))

	var ui = UserInfo{u, "2"}
	strU1, _ := json.Marshal(ui)
	fmt.Println(string(strU1))

	var um = UserMore{ui,"3"}
	strUm, _ := json.Marshal(um)
	fmt.Println(string(strUm))
}

```

- struct中如何屏蔽空值字段
设置 omitempty                   
```go
type User struct {
	Name     string `json:"name"`
	Pwd      string `json:"pwd"`
	Info     string `json:"info, omitempty"`
	MoreInfo string `json:"more_info"`
}
```

- client向server发送10000个数据，server一次只能处理100个。                     
chan即可。                     

- 项目中遇到的坑                   
    - grpc遇到的坑     
        - 报status code=unavailable
            - 因网络抖动，导致没有连接上。                                 
            - 实例有问题，无法处理客户端请求。   
        - 解决：
            - 增加重连机制。               
                - 代码端增加重连
            - 重试机制。
                - retryPolicy中设置retryableStatusCode：针对某种情况进行重连 2017年提案。               
                
    - 内存溢出定位
        - goroutine泄漏问题排查
            我们一个内部小服务需要ssh到目标环境上执行些操作。                                       
            平稳运行了几个月突然就OOM。好在启用了net/http/pprof，调接口查看了goroutine详情等信息。                
            curl http://service:port/debug/pprof/goroutine?debug=1                          
            发现有大量goroutine处于同一状态，应该是在等待数据。有的已经阻塞了几个月...                       
            根据调用栈找到了对应代码的位置，从ssh.dial开始一直到某个地方进行io.readfull阻塞住了。                    
            检查了代码中调用的方法，也设置了超时时间，不会阻塞。研究了下readfull的具体实现逻辑。              
            看逻辑是在对端发送完自己的版本信息后，等待对端的回复，一直没收到消息。为什么没收到消息，有点蒙蔽。
            我们在机器用netstat查看本地建立的连接，有上万个establish状态。再到连接到的机器上检查，发现机器上几乎没有几个连接。                   
            这属于TCP半打开状态，连接成功后，可能对端挂掉了而连接没启动keepalive，导致另一端无法发现这种情况。                 
            在机器上执行了下ss -aeon|grep :36000|grep -v time|wc -l，发现确实没开keeplive。那我们就去开启呗。                
            看了下使用的版本是1.9.2编译的，ssh.Dial用的是net.DialTimeOut方法，返回的net.Conn结构体确实keepalive是默认关闭的。                         
            于是将代码迁移到了新的版本1.14中，检查了net.Conn结构体，返回是开启的。于是打包编译上线。                  
            
            一天后观察，又出现了问题。               
            还是TCP建立了连接，对端不响应。再次检查是否在哪里遗漏了timeout，发现是在handshake的时候，没有作为超时控制的参数使用。                
            而net.Conn的IO等待是非阻塞的epoll_pwait实现的。进入等待的goroutine就会被挂起，直到有事件进来。                  
            于是，我们在ssh.Dial处加个下net的setDeadLine()方法，设置了超时。编译上线后，恢复正常。             
          
            这个只解决了出现异常时及时关闭连接，没有解决可能造成异常的情况。不过如果虚机都异常了，再跑任务也没什么意义。


- 如何实现跨域请求
    - nginx使用反向代理功能，将本地URL映射到要跨域访问的服务器上。nginx通过检查url前缀，把HTTP请求转发到后面真实的服务器上。
      并通过rewrite把前缀再去掉。                                                   
      将nginx.conf文件的server节点文件的参数改成如下内容。                
      
```shell
server {
        location / {
            root   html;
            index  index.html index.htm;
            //允许cros跨域访问
            add_header 'Access-Control-Allow-Origin' '*';

        }
        //自定义本地路径
        location /apis {
            rewrite  ^.+apis/?(.*)$ /$1 break;
            include  uwsgi_params;
            proxy_pass   http://www.binghe.com;
       }
}
```

- 服务器端增加cors代码，设置header，access-control-allow-origin,"*"
        

- TCP、UDP的区别
    
    |类别|面向连接|可靠性|速度|数据|用途|
    |:---:|:---:|:---:|:---:|:---:|:---:|
    |TCP|是|可靠|慢|面向字节流|文件传输、邮件等|
    |UDP|否|不可靠|快|面向数据报|视频、音频等|

- 连接是什么意思
    两个PC机上的两个进程，通过端口逻辑建立了通道。
    
- 七层协议的功能介绍一下
   - 应用层：通过应用进程间的交互，完成特定的网络应用。
   - 运输层：负责向两台主机进程之间的通信提供通用的数据传输服务。应用进程通过利用该服务传送应用层的报文。 
   - 网络层：在网络中通信的两个计算机之间可能会经过很多数据链路，网络层就是选择合适的网间路由和交换节点，确保数据传送。
   - 数据链路层：两台主机之间的数据传输，总是在一段一段的链路上传送的，这就需要使用专门的链路层协议。在两个相邻节点之间传送数据时，
               数据链路层将网络层交下来的IP数据报组装称帧，在两个相邻节点间的链路上传送帧。每一帧包括数据和必要的控制信息。
   - 物理层：在物理层传送的数据单位是比特，物理层作用是实现相邻计算机节点之间的比特流的透明传送，尽可能屏蔽掉传输介质和物理设备的差异。

- IP是面向连接的协议吗？
    IP是面向无连接的，仅仅负责将数据传递给目标地址，本身并不保持连接状态。                            

- DOS攻击针对的是TCP三次握手中的哪次？
    攻击者发送一个SYN报文段，当服务器返回ACK后，该攻击者就不再对其进行再确认，那么这个连接就处于半连接状态。服务器就会重复发送ACK给发送方，浪费服务器资源。                            
  
- 为什么挥手需要四次
    服务端收到关闭连接后，需要等server端所有报文都发送完，才能发送FIN报文，因此不能一起发送，所以需要四次。                    
  
- 从技术角度讲下IP V6比IP V4好的原因。
    - 地址长度：V4是32位（4字节）地址长度；V6具有128位（16字节）地址长度。      
    - 地址表示方法：V4是用小数表示的二进制数，V6是用16进制表示的二进制数。
    - 地址配置：V4地址可以手动或者DHCP配置；V6需要ICMP或DHCPV6的无状态地址自动配置（SLAAC）。  
    - 数据包区别：V4 576字节，碎片可选。V6 1280字节，不会碎片。                   
    - 身份验证和加密：V6提供身份验证和加密；V4没有。                 
    

- GMP模型
- GC
- Redis内存淘汰策略
  
- goroutine快的原因
    原来需要进行线程级别的切换，现在线程保持，只需要切换goroutine即可。

- 进程、线程、协程区别

- HTTP状态码401
    未授权 请求要求用户的身份验证。                
  
- 无锁化编程有哪些常见的方法及原理
    - 原子操作
    - 只有一个生产者和消费者，可以做到免锁访问环形缓冲区
    - CAS（compare and swap）无锁队列等待。
    - RCU(Read-Only-Update) 新旧副本切换，对旧副本延迟释放。                        

- Go是面向对象的吗？与Java面向对象设计上有什么不同？
  go既是又不是面向对象。它允许以面向对象风格编程。               
  - Go中没有object概念，只有struct。
  - 面向对象的编程，没有继承。                       
    


        
    

