package db

type Page struct {
	Page   int   `json:"page" form:"page"`
	Size   int   `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result any   `json:"result" form:"result"`
	Last   int64 `json:"last" form:"last"`
}
type PageT[T any] struct {
	Page   int   `json:"page" form:"page"`
	Size   int   `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result []*T  `json:"result" form:"result"`
	Last   int64 `json:"last" form:"last"`
}

// ToPageT Page ➜ PageT[any]
func (p *Page) ToPageT() *PageT[any] {
	result, ok := p.Result.([]*any)
	if !ok {
		// 如果不是 []*any，尝试强制转换 []interface{} -> []*any
		if slice, ok := p.Result.([]interface{}); ok {
			var casted []*any
			for _, v := range slice {
				val := v
				casted = append(casted, &val)
			}
			return &PageT[any]{
				Page:   p.Page,
				Size:   p.Size,
				Total:  p.Total,
				Result: casted,
				Last:   p.Last,
			}
		}
		// 无法转换，返回空结果
		return &PageT[any]{
			Page:   p.Page,
			Size:   p.Size,
			Total:  p.Total,
			Result: []*any{},
			Last:   p.Last,
		}
	}

	return &PageT[any]{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: result,
		Last:   p.Last,
	}
}

// ToPage PageT[T] ➜ Page（可用于 JSON 输出、统一处理等）
func (p PageT[T]) ToPage() *Page {
	// Page.Result 使用 any 类型表示，可以直接赋值
	return &Page{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: p.Result,
		Last:   p.Last,
	}
}
