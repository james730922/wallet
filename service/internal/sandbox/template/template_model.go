package template

import "fmt"

type ITemplate interface {
	Do()
}

type template struct {
}

func NewTemplate() ITemplate {
	return &template{}
}

func (t *template) Do() {
	fmt.Println("It's a Template!")
}
