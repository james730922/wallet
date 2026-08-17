package geetestsdk

import "fmt"

/**
 * sdk lib包的返回结果信息。
 *
 * @author liuquan@geetest.com
 */
type LibResult struct {
	Status int
	Data   string
	Msg    string
}

func NewGeetestLibResult() *LibResult {
	return &LibResult{0, "", ""}
}

func (g *LibResult) setAll(status int, data string, msg string) {
	g.Status = status
	g.Data = data
	g.Msg = msg
}

func (g *LibResult) String() string {
	return fmt.Sprintf("LibResult{Status=%d, Data=%s, Msg=%s}", g.Status, g.Data, g.Msg)
}
