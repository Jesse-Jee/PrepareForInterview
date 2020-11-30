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
    示例：     
```go
    type SimpleFactoryOfParse interface {
    	Parse(name string)
    }
    
    type AConfig struct {
    }
    
    type BConfig struct {
    }
    
    func (a AConfig) Parse(name string) {
    	// TODO 解析逻辑
    }
    
    func (b BConfig) Parse(name string) {
    	// TODO 解析实际逻辑
    }
    
    func NewConfigParse(name string) SimpleFactoryOfParse {
    	switch name {
    	case "json":
    		return AConfig{}
    	case "xml":
    		return BConfig{}
    	}
    	return nil
    }
    
```
- 工厂方法
    - 当对象的创建逻辑比较复杂，不是简单new一下，而要组合其他对象，做各种初始化时，使用工厂方法模式。     
        示例： 
```go
    type FactoryMethodOfParse interface {
    	CreateParser() SimpleFactoryOfParse
    }
    
    type AConfigParseFactory struct {
    }
    
    func (a AConfigParseFactory) CreateParser() SimpleFactoryOfParse {
    	return AConfig{}
    }
    
    type BConfigParseFactory struct {
    }
    
    func (b BConfigParseFactory) CreateParser() SimpleFactoryOfParse {
    	return BConfig{}
    }
    
    func NewConfigParserFactory(name string) FactoryMethodOfParse {
    	switch name {
    	case "json":
    		return AConfigParseFactory{}
    	case "xml":
    		return BConfigParseFactory{}
    	}
    	return nil
    }
```    

    
- 抽象工厂
    - 为创建一组或相互依赖的对象提供一个接口，而且无需指定他们的具体类。     
        示例：     
        抽象工厂：一个车库概率
        具体工厂： 一个具体的玛莎车库，奔驰车库，奥迪车库
        抽象产品： 车的概念，包的概念
        具体产品： 奥迪，阿迪达斯
        


# 责任链模式

# 装饰器模式


