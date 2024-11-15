package go_logger

type Option interface {
	Name() string
	Value() interface{}
}
