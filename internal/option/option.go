package option

type Interface interface {
	Name() string
	Value() interface{}
}

type Option struct {
	name  string      `json:"name"`
	value interface{} `json:"value"`
}

func (o *Option) Name() string {
	return o.name
}

func (o *Option) Value() interface{} {
	return o.value
}

func NewOption(name string, value interface{}) *Option {
	return &Option{name: name, value: value}

}
