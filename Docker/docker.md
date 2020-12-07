# Docker如何管理CPU和内存资源

# docker隔离是怎么实现的。
## 容器vs虚拟机
都提供了隔离应用和依赖环境的能力。           
都可以看做是一个沙箱环境，使应用可以部署在任意宿主机。                 
不同的是，虚拟机需要依赖硬件设备来提供资源隔离。            

### 虚拟机
虚拟机是一个真实计算机操作系统的封装，它运行在物理设备上，通过hypervisor进行建立和运行虚拟机体系。              

|vm|                
|:---:|          
|hypervisor|
|host os|
|server|



VM：             

|app|
|:---:|
|bins/libs|
|guest os|

在host os基础上，通过hypervisor进行虚拟资源控制，并拥有自己的guest os, 虽然隔离更彻底，但资源开销更大。                   


### 容器
容器提供的是操作系统级别的进程隔离。docker本身只是操作系统的一个进程。在容器技术下，进程之间的网络和空间等都是隔离的。              
多个容器之间共享了宿主机的操作系统内核，在host OS上通过docker engine共享host OS的内核。           



## docker资源隔离namespace
docker利用Linux namespace来实现多个容器的互相隔离。            
- mount namespace 用于隔离文件系统的挂载点
- UTS namespace 用于隔离hostname和domainname
- IPC namespace 用于隔离进程间通信
- PID namespace 用于隔离进程ID
- network namespace 用于隔离网络
- user namespace 用于隔离用户和用户组                 





    
              



              
