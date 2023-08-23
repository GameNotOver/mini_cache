package di

import "go.uber.org/dig"

type ContainerInterface interface {
	RegisterWithName(constructor interface{}, name string) error
	Register(constructor interface{}, opts ...dig.ProvideOption) error
	MustRegister(constructor interface{}, opts ...dig.ProvideOption)
	Call(function interface{}, opts ...dig.InvokeOption) error
	MustCall(function interface{}, opts ...dig.InvokeOption)
}

// Global 提供全局注册能力.
var Global = New()

func New() *Container {
	return &Container{dig.New()}
}

type Container struct {
	*dig.Container
}

func (c *Container) Register(constructor interface{}, opts ...dig.ProvideOption) error {
	return c.Provide(constructor, opts...)
}

func (c *Container) RegisterWithName(constructor interface{}, name string) error {
	return c.Register(constructor, dig.Name(name))
}

func (c *Container) Call(function interface{}, opts ...dig.InvokeOption) error {
	return c.Invoke(function, opts...)
}

func (c *Container) MustRegister(constructor interface{}, opts ...dig.ProvideOption) {
	if err := c.Register(constructor, opts...); err != nil {
		panic(err)
	}
}

func (c *Container) MustCall(function interface{}, opts ...dig.InvokeOption) {
	if err := c.Call(function, opts...); err != nil {
		panic(err)
	}
}

func MustRegister(constructor interface{}, opts ...dig.ProvideOption) {
	Global.MustRegister(constructor, opts...)
}
