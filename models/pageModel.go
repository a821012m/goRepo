package models

type PageModel struct {
	Index    int64 `uri:"index" binding:"required,min=1"`
	PageSize int64 `uri:"pageSize" binding:"required,min=1"`
}

type PageResult struct {
	PageInfo   PageModel
	TotalCount int64
	TotalPage  int64
	Datas      interface{}
}
