# 延时请求
因未满足条件而暂时无法被处理的kafka请求。                           
延时请求处理机制：分层时间轮算法                    

时间轮：像是手表，时针，分针，秒针。各自有自己的刻度，又彼此相关。秒针转一圈，分针转一格。               
这就是典型的分层时间轮。                    

kafka对应手表的一格，叫一个 bucket                 
一个推进，叫做一个 滴答  tick                  

bucket 底层是双向循环链表。插入和删除时间复杂度是O(1)。               


## 源码实现
TimingWheel->N个TimerTaskList（Bucket上）-> N个TimerTaskEntry -> TimerTask

timerTaskEntry与TimerTask 一对一。           


**TimerTask**               
每个timerTask都有一个delayMs超时时间字段。               
这个类绑定了一个TimerTaskEntry字段，因为每个定时任务需要知道自己在哪个bucket链表下的哪个链表元素上。            
 
在往timerTaskEntry赋值之前，需先考虑这个定时任务是否已经绑定了其他的timerTaskEntry，如果是就必须先取消绑定。                

**TimerTaskEntry**                       
timerTaskEntry表征的是bucket链表下的一个元素。                   
它定义了timeTask类字段，用来指定定时任务。                       
封装了过期时间戳字段，定义了定时任务的过期时间。                    
定义了list,prev,next，分别表示bucket实例，自身的prev，next指针。                  

**如何删除一个定时任务？**
调用TimerTaskList中的remove方法，remove就是将任务从双向循环链表中移除。置空timerTaskEntry的list，如果这个为空，那么这个实例timerTaskEntry就变成了孤儿。               


**TimerTaskList**           
定义了root节点，同时定义了：                    
- taskCounter 标识当前链表中的总定时任务数。               
- expiration 表示链表所在bucket的过期时间戳                     


timerTaskList提供了add和remove方法，实现了将定时任务插入到链表和从链表中移除定时任务的逻辑。                       
flush方法是清空链表中的所有元素，并对每个元素执行指定的逻辑。                 


**TimerWheel**              
- tickMs:滴答一次时长。第一层时间轮的时间是1毫秒
- wheelSize: 每一层时间轮上的bucket数量，第一层bucket数量是20            
- startMs: 时间轮对象被创建的起始时间戳
- taskCounter: 这一层时间轮上的总定时任务数。                  
- queue: 将所有bucket按照过期时间排序的延迟队列。                    
- interval: 这一层时间轮总时长。 相当于滴答时长乘以数量                  
- buckets: 时间轮上所有timerTaskList对象。                   
- currentTime: 当前时间戳                    
- overflowWheel: 按需创建时间轮。尝试放到第一层时间轮，第一层放不下了，尝试创建第二层时间轮，并再次尝试放入。以此类推。                

上层所用的滴答时长等于下层时间轮总时长。


**add方法**
新增一个定时任务逻辑

- 1 获取定时任务时间戳
- 2 查看定时任务是否已取消             
- 3 计算目标bucket序号，看是保存在哪个bucket中。                   
- 4 如果这个bucket是首次插入定时任务，那么这个bucket还要加入到delayqueue中，方便kafka轻松获取那些已过期的bucket，并删除它们。                 

**advanceClock**            
向前驱动时钟。             
timeMs表示要把时钟向前推动到这个时点。向前驱动的时点必须要超过bucket的时间范围，才是有意义的推进。                     
如果超出了bucket的时间范围，代码就会更新当前时间到下一个bucket的起始时点，并递归的为上一层时间轮做向前推进动作。                  
                 














































 
 




                          






