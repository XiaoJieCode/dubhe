package db

type Page struct {
	Page   int   `json:"page" form:"page"`
	Size   int   `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result []any `json:"result" form:"result"`
	Last   int64 `json:"last" form:"last"`
}

type PageT[T any] struct {
	Page   int   `json:"page" form:"page"`
	Size   int   `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result []T   `json:"result" form:"result"`
	Last   int64 `json:"last" form:"last"`
}

// ToPageT 将非泛型分页转换成泛型分页，适用于泛型处理
func (p *Page) ToPageT() *PageT[any] {
	// 如果不能断言为 []any，返回空结果
	return &PageT[any]{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: []any{},
		Last:   p.Last,
	}
}

// ToPage 将泛型分页转换为非泛型分页，方便统一JSON输出或兼容旧接口
func (p PageT[T]) ToPage() *Page {
	return &Page{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: any(p.Result).([]any),
		Last:   p.Last,
	}
}

// ToAnySlice 辅助函数，将任意类型切片转换成 []any，方便赋值给Page.Result
func ToAnySlice[T any](src []T) []any {
	res := make([]any, len(src))
	for i, v := range src {
		res[i] = v
	}
	return res
}
