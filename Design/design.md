# 单例模式
## 定义
一个类只允许创建一个对象。  
## 使用
- 表示全局唯一类： config配置
- 处理资源访问冲突：项目logger

## 分类
- 饿汉式
    - 类加载的时候就给初始化好
- 懒汉
    - 实际要用到类的时候才给初始化
- 双重检测

## 实践
sync.Once 就是一个单例模式。     
在实际使用中，可以用来初始化全局唯一类之类的。     
如：
```go
    
type Config struct {
	Port int
}

var conf *Config

func GetConfig() *Config {
	if conf == nil {
		conf = &Config{Port: 8080}
	}
	return conf
}

var once sync.Once
// once的使用
// 实际上go自身实现了一个单例模式的once，多用于初次加载配置。
func GetconfigOnce() *Config {
	once.Do(func() {
		conf = &Config{Port: 8080}
	})
	return conf
}

func main() {
	c1 := GetConfig()
	c2 := GetConfig()
	fmt.Println(c1 == c2)
	// c1,c2取到的是同一个对象。
}
```

# 工厂模式
一般分为三类：
- 简单工厂模式
    - 创建型模式，由一个工厂对象决定创建出哪种产品类型的实例。      
- 工厂方法
- 抽象工厂


