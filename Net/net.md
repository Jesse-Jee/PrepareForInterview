# NAT  
## 介绍  
NAT是network address translation的缩写，网络地址转换的缩写。
因为IP地址有限，所以需要多人共用一个公网IP来访问互联网。当需要访问互联网时，网络边界设备（路由器，防火墙等）将各个主机的私网IP地址，转成公网IP地址。
即，将数据包中的IP地址替换为其他IP地址的功能。
## 实现方式  
### 静态NAT  
一个公网IP对应一个私网IP，一对一转换。
### 动态NAT  
在路由器上配置一个公网IP池，当私有IP要进行外部通信时，从池中选择一个公网IP，并将他们的关系绑定到NAT表中。通信结束后，这个公网IP才会被释放，供其他私有IP使用。

### NAPT   
端口地址复用技术。是最常用的，提供一对多的方式。对外只有一个公网IP地址，内部通过端口来区分不同的私网IP主机数据。

## NAT如何区分不同私网的host呢？ 
对于TCP/UDP使用
私网host + port <---> 公网host + port 

对于ICMP使用
私网host + sessionID <----> 公网host + sessionID

对于其他协议，NAT使用的也是类似的转换规则，选择将host能轻易区分出来的字段作为key值，动态创建映射表项，做双向地址+key的转换。

# ARP 
为IP地址到对应的硬件地址提供动态映射。

ARP请求（以太网广播） 问 某个地址的mac地址是多少啊？ 这个mac地址的给它单播返回告诉它。


# ICMP 
IP层的一个组成部分。用于传递差错报文及其他需要注意的信息。
ICMP报文：
| 8位类型 ｜ 8位代码 ｜ 16位校验和 ｜  

## Ping 

## traceroute 程序
使用ICMP报文和IP首部的TTL字段（生存周期），TTL字段是由发送端初始设置的一个8bit字段。
ICMP设置的TTL最大是255

每个处理数据报的路由器都需要把TTL的值减1或减去在路由器中停留的秒数。
当路由器收到一份数据报的TTL是0或者1时，则路由器不转发这个数据报。路由器会把其丢弃，并给信源饥发送一个ICMP超时信息。

traceroute发送udp包给目的主机

# TCP
链路层：  
以太网：采用CSMA/CD的媒体接入方法。速率为10Mb/s，mac地址为48bit  
IEEE802：


封装格式：  
802.3:  
| 目的地址 6 ｜ 源地址 6 | 长度 2 | DSAP  1 | SSAP 1 | cntl 1 | org code 3 | 类型 2 ｜ 数据 38-1492 ｜ CRC 4 ｜  

以太网：  
| 目的地址 6 | 源地址 6 | 类型 2 | 数据 46-1500| CRC 4|  

IP协议：  
类型     IP数据报  
0800 

ARP协议：  
｜类型 ｜ ARP请求/应答 ｜ PAD ｜  
｜0806 ｜    26      ｜ 18

以太网头部： 目的地址  源地址 类型  
类型字段中，放协议，如上面的IP协议
以太网要求数据部分最少要求46字节，不足46字节的要用PAD补齐，超出1500字节的，要分页。  

应用层 FTP  21
传输层 TCP UDP
网络层 IP ICMP
数据层 


# SOCKET 
## 网络进程中如何通信
使用TCP/IP协议的应用程序通常使用socket来进行网络通信。
### 如何确定唯一标识？
网络层 IP地址+ 传输层 协议+端口 可以唯一标识主机中的应用程序。

### 什么是socket
socket就像是一种特殊的文件，一些socket函数对其进行操作（读写i/o，打开，关闭）

### socket的基本操作
#### socket()函数  
创建一个socket描述符，唯一标识一个socket。后续操作通过它进行一些读写操作。  
**int socket(int domain, int type, int protocol);**  
创建socket时可以指定不同的参数
- domain：协议域；AF_INET、AF_INET6、AF_LOCAL（或称AF_UNIX，Unix域socket）、AF_ROUTE等等。协议族决定了socket的地址类型，
在通信中必须采用对应的地址，如AF_INET决定了要用ipv4地址（32位的）与端口号（16位的）的组合、AF_UNIX决定了要用一个绝对路径名作为地址。
- type：socket类型。SOCK_STREAM、SOCK_DGRAM、SOCK_RAW、SOCK_PACKET、SOCK_SEQPACKET等等  
- protocol：协议。IPPROTO_TCP、IPPTOTO_UDP、IPPROTO_SCTP、IPPROTO_TIPC等，
它们分别对应TCP传输协议、UDP传输协议、STCP传输协议、TIPC传输协议。  

当我们调用socket()时，返回的socket描述字存在于协议域空间中。如果想给它赋一个地址，就必须调用bind()函数。  

##### socket类型有哪些？
SOCK_STREAM，SOCK_DGRAM，SOCK_RAW，SOCK_PACKET，SOCK_SEQPACKET。


#### bind()函数  
把一个地址族中的地址赋值给socket。   
**int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen);** 
- sockfd 即socket描述字
- sockaddr 指向要绑定给socket的协议地址，指针。根据地址协议族的不同而不同。
- addrlen 地址长度

通常服务器在启动时会绑定一个众所周知的地址（如IP地址+端口号）用于提供服务，客户端可以通过IP地址加端口号连接服务器。而客户端不用指定。
由系统自动分配一个端口号和自身IP地址的组合。这就是为什么通常服务器在listen之前需要调用bind，而客户端不用调用，而是在connect时由系统随机生成一个。

#### listen()函数
调用listen监听这个socket，这时如果客户端调用connect发出连接请求，服务器就会接收到这个请求。
**int listen(int sockfd, int backlog);**
- sockfd 即要监听的socket描述字
- backlog 即相应socket可以排队的最大连接个数。
socket()函数创建的socket默认是一个主动类型的，listen函数将socket变为被动类型的，等待客户的连接请求。  

**int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen);**

- sockfd 即客户端socket描述字
- addr 即服务器socket地址。
- addrlen 即socket地址长度

客户端通过connect函数与TCP服务端建立连接。

#### accept()函数
TCP服务器在经过socket,bind,listen后，监听端口，当监听到连接请求后，调用accept进行连接。之后就可以进行网络i/o操作了。
就像对文件i/o的操作一样。  
**int accept(int sockfd, struct sockaddr *addr, socklen_t *addrlen);**
- sockfd 服务器socket的描述字
- addr 返回客户端的协议地址
- addrlen 协议地址长度


#### read/write 等函数
万事具备，只欠i/o操作。
- read()/write()
- recv()/send()
- readv()/writev()
- recvmsg()/sendmsg()
- recvfrom()/sendto()


       ssize_t read(int fd, void *buf, size_t count);
       ssize_t write(int fd, const void *buf, size_t count);
       ssize_t send(int sockfd, const void *buf, size_t len, int flags);
       ssize_t recv(int sockfd, void *buf, size_t len, int flags);
       ssize_t sendto(int sockfd, const void *buf, size_t len, int flags,
                      const struct sockaddr *dest_addr, socklen_t addrlen);
       ssize_t recvfrom(int sockfd, void *buf, size_t len, int flags,
                        struct sockaddr *src_addr, socklen_t *addrlen);
       ssize_t sendmsg(int sockfd, const struct msghdr *msg, int flags);
       ssize_t recvmsg(int sockfd, struct msghdr *msg, int flags);

read函数是负责从fd中读取内容.当读成功时，read返回实际所读的字节数，如果返回的值是0表示已经读到文件的结束了，小于0表示出现了错误。
如果错误为EINTR说明读是由中断引起的，如果是ECONNREST表示网络连接出了问题。
write函数将buf中的nbytes字节内容写入文件描述符fd.成功时返回写的字节 数。失败时返回-1，并设置errno变量。
在网络程序中，当我们向套接字文件描述符写时有俩种可能。1)write的返回值大于0，表示写了部分或者是 全部的数据。2)返回的值小于0，此时出现了错误。
我们要根据错误类型来处理。如果错误为EINTR表示在写的时候出现了中断错误。如果为EPIPE表示 网络连接出现了问题(对方已经关闭了连接)。



#### close()函数
关闭相应socket描述字
close操作只是使相应socket描述字的引用计数-1，只有当引用计数为0的时候，才会触发TCP客户端向服务器发送终止连接请求。




# UDP 
## 什么是UDP
- UDP是一个简单的面向数据报的运输层协议：进程的每个输出操作，都正好产生一个UDP数据报，并组成一份待发送的IP数据包。
- UDP不提供可靠性：它把应用程序传给IP层的数据发送出去，但是并不能保证它们到达目的地。
- 应用程序必须关心IP数据报的长度，如果它超过网络的MTU，那就要对IP数据报进行分片。

## UDP的封装
![Image text](https://raw.githubusercontent.com/jizengguang/PrepareForInterview/master/Picture/udp.png)

## UDP的三大应用
- 查询类：DNS
  - 没有TCP的三次握手过程，快
  - 多个DNS同时查询
- 数据传输：TFTP
  - 停止等待协议，慢
  - 适合于无盘工作站
- 语音视频流：
  - 支持组播和广播
  - 支持丢包，保障效率
  
## UDP首部
![Image_text](https://raw.githubusercontent.com/jizengguang/PrepareForInterview/master/Picture/udp_header.png)


## UDP校验和
- UDP校验和覆盖UDP首部和UDP数据
- IP校验和只覆盖IP首部
- UDP的校验和是可选的，而TCP的校验和是必需的。

## IP分片
- IP把MTU和数据报长度进行比较。
- 把一份IP数据报分片以后，只有到达目的地才进行重新组装。
- 重新组装由目的端的IP层来完成，其目的是使分片和重组过程对运输层（TCP、UDP）透明
- 已经分片过的数据报有可能会再次进行分片。

- 在分片时，除最后一片外，其他每一片中的数据部分（除IP首部外的其余部分）必须是8字节的整数倍。
- IP首部被复制到各个片中，但是端口号在UDP首部，只能在第一片中被发现。

# TCP UDP区别
类型：面向连接：传输可靠性：传输形式：传输效率：所需资源：应用场景：首部字节  
TCP：面向连接，可靠，字节流，慢，多，文件传输/邮件传输，20-60  
UDP：无连接，不可靠，数据报，快，少，视频，音频，8  


# OSI体系结构
应用层
表示层
会话层
运输层
网络层
数据链路层
物理层


# 五层协议结构
应用层 FTP/DNS/HTTP协议
运输层 TCP/UDP
网络层 IP
数据链路层
物理层

## 各层作用
### 应用层
通过应用进程间的交互，完成特定的网络应用。

### 运输层
负责向两台主机进程之间的通信提供通用的数据传输服务。应用进程通过利用该服务传送应用层的报文。

### 网络层
在网络中通信的两个计算机之间可能会经过很多数据链路，网络层就是选择合适的网间路由和交换节点，确保数据传送。

### 数据链路层
两台主机之间的数据传输，总是在一段一段的链路上传送的，这就需要使用专门的链路层协议。在两个相邻节点之间传送数据时，
数据链路层将网络层交下来的IP数据报组装称帧，在两个相邻节点间的链路上传送帧。每一帧包括数据和必要的控制信息。

### 物理层
在物理层传送的数据单位是比特，物理层作用是实现相邻计算机节点之间的比特流的透明传送，尽可能屏蔽掉传输介质和物理设备的差异。


## DNS
域名系统，将域名与IP地址相互映射。
## HTTP
超文本传输协议

# 三次握手
客户端： 发送带SYN标志的数据包，例如SYN J 第一次握手
服务端： 接收到SYN标志的数据包，返回给客户端 SYN K , ACK seq J+1 ，带ACK意思是这是一个应答包
客户端： 接收到服务端返回的SYN/ACK数据包，返回给服务端 ACK K+1

![Image_text](https://raw.githubusercontent.com/jizengguang/PrepareForInterview/master/Picture/three_hand.png)

## SYN和ACK
SYN是同步序列号，是TCP/IP建立连接时使用的握手信号。
ACK是确认字符，在数据通信中，接收站发给发送站的一种传输类控制字符。表示发来的数据已确认接收无误。

## 为什么需要三次握手？
保证数据的正常发送与接收
第一次握手，客户端无法确认；server端确认对方发送正常，自己接收正常。
第二次握手，客户端确认自己发送，接收正常，对方发送接收正常；server端确认对方发送正常，自己接收正常。
第三次握手，客户端确认自己发送，接收正常，对方发送接收正常；server端确认对方发送接收正常，自己发送接收正常。

# 四次挥手 
- 客户端进程发出连接释放报文，并且停止发送数据。释放数据报文首部，FIN=1，其序列号为seq=u（等于前面已经传送过来的数据的最后一个字节的序号加1），
此时，客户端进入FIN-WAIT-1（终止等待1）状态。 TCP规定，FIN报文段即使不携带数据，也要消耗一个序号。
- 服务器收到连接释放报文，发出确认报文，ACK=1，ack=u+1，并且带上自己的序列号seq=v，此时，服务端就进入了CLOSE-WAIT（关闭等待）状态。
TCP服务器通知高层的应用进程，客户端向服务器的方向就释放了，这时候处于半关闭状态，即客户端已经没有数据要发送了，但是服务器若发送数据，
客户端依然要接受。这个状态还要持续一段时间，也就是整个CLOSE-WAIT状态持续的时间。
- 客户端收到服务器的确认请求后，此时，客户端就进入FIN-WAIT-2（终止等待2）状态，等待服务器发送连接释放报文（在这之前还需要接受服务器发送的最后的数据）。
- 服务器将最后的数据发送完毕后，就向客户端发送连接释放报文，FIN=1，ack=u+1，由于在半关闭状态，服务器很可能又发送了一些数据，
假定此时的序列号为seq=w，此时，服务器就进入了LAST-ACK（最后确认）状态，等待客户端的确认。
- 客户端收到服务器的连接释放报文后，必须发出确认，ACK=1，ack=w+1，而自己的序列号是seq=u+1，此时，客户端就进入了TIME-WAIT（时间等待）状态。
注意此时TCP连接还没有释放，必须经过2MSL（最长报文段寿命）的时间后，当客户端撤销相应的TCB后，才进入CLOSED状态。
- 服务器只要收到了客户端发出的确认，立即进入CLOSED状态。同样，撤销TCB后，就结束了这次的TCP连接。可以看到，服务器结束TCP连接的时间要比客户端早一些。
![Image_text](https://raw.githubusercontent.com/jizengguang/PrepareForInterview/master/Picture/four_hands.png)

## 四次挥手中，服务器端的Close-Wait有什么用？
用于将最后的数据发送完毕，再向客户端发送连接释放报文。

## 四次挥手中，为什么客户端需要等待2MSL？
第一，保证客户端发送的最后一个ACK报文能够到达服务器，因为这个ACK报文可能丢失，站在服务器的角度看来，我已经发送了FIN+ACK报文请求断开了，
客户端还没有给我回应，应该是我发送的请求断开报文它没有收到，于是服务器又会重新发送一次，而客户端就能在这个2MSL时间段内收到这个重传的报文，
接着给出回应报文，并且会重启2MSL计时器。
第二，防止类似与“三次握手”中提到了的“已经失效的连接请求报文段”出现在本连接中。客户端发送完最后一个确认报文后，在这个2MSL时间中，
就可以使本连接持续的时间内所产生的所有报文段都从网络中消失。这样新的连接中不会出现旧连接的请求报文。
 
 
# TCP协议如何保证可靠传输
1. 应用数据被分割称TCP认为合适发送的数据块。
2. TCP会给发送的每个包进行编号，接收方对数据包进行排序，把有序数据传送给应用层。
3. 校验和：TCP将保持它首部和数据的校验和，目的是检测数据在传输过程中的变化。如果收到的校验和有差错，将丢弃此报文段和不确认收到此报文段。
4. TCP的接收端会丢弃重复的数据。
5. 流量控制： TCP连接的每一方都有固定大小的缓冲空间，TCP的接收端只允许发送端发送接收端缓冲区能接纳的数据。当接收方来不及处理发送方的数据，
能提示发送方降低发送速率，防止包丢失。TCP使用的流量控制协议是可变大小的滑动窗口协议。
6. 拥塞控制：当网络拥塞时，减少数据发送。
7. ARQ协议：每发完一个分组就停止发送，等待对方确认，在收到确认后再发下一个分组。
8. 超时重传：当TCP发出一个段后，它启动一个定时器，等待目的端确认收到这个报文段。如果不能及时收到一个确认，将重发这个报文段。


# TCP keep-alive机制

# HTTP1 2 3区别

# HTTPS


# TCP滑动窗口，ACK机制

# tcp报文结构


# HTTP重定向机制

# TCP粘包


# Session，Cookie区别


# 网络协议栈


# HTTPS工作过程


# HTTP请求报文格式



# 一次URL发生了什么



# 对称加密和非对称加密


# tcp能建立多少链接


# post get 区别


# DNS协议解析过程





