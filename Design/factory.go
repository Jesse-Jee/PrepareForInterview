package Design

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
