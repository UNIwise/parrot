package poedit

type PoeditClient interface{}

type PoeditClientImpl struct{}

func NewClient() *PoeditClientImpl {
	return &PoeditClientImpl{}
}
