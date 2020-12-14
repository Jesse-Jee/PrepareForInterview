# gin
一个web框架，实际上就是封装了net/http的web框架。
增加了中间件支持。使用前缀树路由。

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
添加URI时，会找到方法对应的method trees添加进去。如果树是空的就初始化一个树。
addRoute不是并发安全的。        

node中包含了路径信息，它的类型信息等等。                      

````go
const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type node struct {
	path      string
	indices   string
	wildChild bool
	nType     nodeType
	priority  uint32
	children  []*node
	handlers  HandlersChain
	fullPath  string
}
````




## radix-tree举例
1 romane                
2 romanus               
3 romulus               
4 rubens                
5 ruber                 
6 rubicon                   
7 rubicundus                    

![Image_text](https://raw.githubusercontent.com/Jesse-Jee/PrepareForInterview/master/Picture/radix-tree.png)




# group
是一个RouterGroup。
gin可以通过engine.Group设置路由组。可以添加公共中间件，或添加相同前缀的路由。
```go
    func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
    	return &RouterGroup{
    		Handlers: group.combineHandlers(handlers),
    		basePath: group.calculateAbsolutePath(relativePath),
    		engine:   group.engine,
    	}
    }
```
可以使用group设置访问不同版本的接口。

```go
    func main() {
    	router := gin.Default()
    
    	// Simple group: v1
    	v1 := router.Group("/v1")
    	{
    		v1.POST("/login", loginEndpoint)
    		v1.POST("/submit", submitEndpoint)
    		v1.POST("/read", readEndpoint)
    	}
    
    	// Simple group: v2
    	v2 := router.Group("/v2")
    	{
    		v2.POST("/login", loginEndpoint)
    		v2.POST("/submit", submitEndpoint)
    		v2.POST("/read", readEndpoint)
    	}
    
    	router.Run(":8080")
    }
```

# gin.Default() 和 gin.New()的区别
default默认包含了logger和recovery中间件。logger日志中间件，Recovery()从任何panic中recover恢复，返回500的中间件。
new默认不包含任何中间件。想要使用自己的中间件，请调用engine.Use(mid).

# 使用中间件
Use()函数
```go
    // Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
    // included in the handlers chain for every single request. Even 404, 405, static files...
    // For example, this is the right place for a logger or error management middleware.
    func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
    	engine.RouterGroup.Use(middleware...)
    	engine.rebuild404Handlers()
    	engine.rebuild405Handlers()
    	return engine
    }
```
将一个全局中间件连接到路由上。中间件将包含在请求的处理链中。

gin的中间件本质为handlerFunc,只要实现一个handlerFunc就可以自定义一个中间件。         
gin中间件实现原理是设计模式中的责任链模式，责任链模式就是为请求创建一个对象链，对象链上的每个对象都可以依次对请求进行处理，并把处理过的请求
传递给下一个对象。

# 设计模式
责任链模式。
