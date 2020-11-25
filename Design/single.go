package Design

import (
	"fmt"
	"sync"
)

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
