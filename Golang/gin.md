# gin
一个web框架

## 特性
- 快，基于radix树的路由，小内存占用，性能好。
- 支持中间件，logger，鉴权，db等。
- 支持crash，可以catch到HTTP请求中的panic并recover它。
- 可以解析json
- 支持路由组
- 错误处理可记录
- 内置渲染


## 路由
基于 radix树实现。            
不同的method（get/post/delete/put等）下面挂不同的路由树。  
gin.Default()返回一个engine对象，engine对象继承自routergroup。  
默认初始化时，routergroup设置为：          
```go
    RouterGroup: RouterGroup{
    			Handlers: nil,
    			basePath: "/",
    			root:     true,
    		},
```


handle方法，传入HTTP方法，相对路径，handles调用链。      
添加URI时，会找到方法对应的method tress添加进去。如果树是空的就初始化一个树。
addRoute不是并发安全的。        
      






# 中间件

# group
是一个RouterGroup

# 设计模式
