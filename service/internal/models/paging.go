package models

import "github.com/james730922/wallet/service/internal/pb/zqbapis"

const (
	pagingDefaultIndex = 1
	pagingDefaultSize  = 25
)

func NewPaging(pi, ps int) Paging {
	p := Paging{
		Index: pi,
		Size:  ps,
	}
	return p
}

func NewPagingFromProto(pp *zqbapis.Paging) Paging {
	p := Paging{}

	if pp != nil {
		p.Index = int(pp.PageIndex)
		p.Size = int(pp.PageSize)
	}

	return p
}

type Paging struct {
	Index int `form:"pi"` // 頁碼
	Size  int `form:"ps"` // 比數
}

func (p *Paging) GetIndex() int {
	if p.Index < 1 {
		return pagingDefaultIndex
	}
	return p.Index
}

func (p *Paging) GetSize() int {
	if p.Size < 1 {
		return pagingDefaultSize
	}
	return p.Size
}

func (p *Paging) GetOffset() int {
	return p.GetSize() * (p.GetIndex() - 1)
}

func (p *Paging) GetPaging() (int, int) {
	return p.GetIndex(), p.GetSize()
}

func (p *Paging) PagingPtr() *Paging {
	return p
}

func NewPagingResult(paging *Paging, count int) *PagingResult {
	totalPage := count / paging.GetSize()
	if count%paging.GetSize() > 0 {
		totalPage++
	}

	return &PagingResult{
		Index:     paging.GetIndex(),
		Size:      paging.GetSize(),
		TotalPage: totalPage,
		TotalRow:  count,
	}
}

type PagingResult struct {
	Index     int `json:"pi"`         // 頁碼
	Size      int `json:"ps"`         // 比數
	TotalPage int `json:"total_page"` // 總頁數
	TotalRow  int `json:"total_row"`  // 總筆數
}
