package Design

// 简单工厂模式
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

// 工厂方法模式

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
