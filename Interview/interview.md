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





    
    
        
    

