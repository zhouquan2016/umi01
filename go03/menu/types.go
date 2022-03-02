package menu

type TreeNode struct {
	Id     uint   `json:"id"`
	PId    string `json:"pId"`
	Value  uint   `json:"value"`
	Title  string `json:"title"`
	IsLeaf bool   `json:"isLeaf"`
}

type AddQuery struct {
	IsLeaf    *bool   `json:"isLeaf"`
	Name      *string `json:"name"`
	ParentId  *uint   `json:"parentId"`
	Path      *string `json:"path"`
	Resources []struct {
		Id   uint   `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
}

type UpdateQuery struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Resources []struct {
		Id   uint   `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
}

type AddResourceQuery struct {
	MenuId uint     `json:"menuId"`
	Codes  []string `json:"codes"`
}
