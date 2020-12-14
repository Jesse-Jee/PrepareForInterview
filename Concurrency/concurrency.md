# master-worker模式
多线程并行计算的一种实现方式。                 
它的核心思想是：启动两个进程协同工作，master进程和worker进程。               
master负责任务的接收和分配，worker负责具体子任务执行。                   
每个worker执行完任务后，把结果返回给master，最后由master汇总结果。                  
分治的思想。                  



