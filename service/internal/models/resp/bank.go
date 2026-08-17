package resp

type BankCodeResp struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status int    `json:"status,string"`
}
