package poedit

type Client interface{}

type ClientImpl struct{}

func NewClient() *ClientImpl {
	return &ClientImpl{}
}
