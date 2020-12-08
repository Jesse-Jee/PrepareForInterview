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
    
    
    
        
    

