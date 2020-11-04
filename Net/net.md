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
网络层IP地址+传输层 协议+端口 可以唯一标识主机中的应用程序。

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




# TCP UDP区别

# TCP keep-alive机制

# HTTP1 2 3区别

# HTTPS

# 三次握手

# 四次挥手


# TCP滑动窗口，ACK机制

# tcp报文结构


# HTTP重定向机制

# TCP粘包


# Session，Cookie区别


# 网络协议栈


# HTTPS工作过程


# HTTP请求报文格式



# 一次URL发生了什么


# close-wait作用


# 对称加密和非对称加密


# osi七层模型


# tcp能建立多少链接


# post get 区别


# DNS协议解析过程





