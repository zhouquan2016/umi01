package conf

type BaseQuery struct {
	SortFiled string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
	PageSize  int    `json:"pageSize"`
	Current   int    `json:"current"`
}

func (query *BaseQuery) Offset() int {
	return query.Current*query.PageSize - query.PageSize - 1
}

type PageResult struct {
	ApiResult
	Current int   `json:"current"`
	Total   int64 `json:"total"`
}

func NewPageResult(current int, total int64, rows interface{}) *PageResult {
	return &PageResult{ApiResult: *SuccessResult(rows), Current: current, Total: total}
}
