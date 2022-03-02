package db

type Identify interface {
	GetId() uint
}
type User struct {
	Id uint `gorm:"primaryKey"`
	//姓名
	Name string
	//密码
	Password string
	//头像地址
	Avatar string
	//邮箱
	Email string
	//个性签名
	Signature string
	//头衔
	Title string
	//地址
	Address string
	//手机号
	Phone string
	//访问权限
	RoleNo string
	RoleId uint
	//系统预留无法删除
	IsSysDefault bool
}

type Menu struct {
	Id           uint        `gorm:"primaryKey" json:"id"`
	Name         string      `json:"name"`
	Children     []*Menu     `json:"children" gorm:"-"`
	IsLeaf       bool        `json:"isLeaf"`
	Path         string      `json:"path"`
	ParentId     uint        `json:"parentId"`
	IsSysDefault bool        `json:"isSysDefault"`
	Dept         string      `json:"dept"`
	DeptName     string      `json:"deptName"`
	Resources    []*Resource `json:"resources" gorm:"-"`
}

type Role struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
	No   string `json:"no"`
	Name string `json:"name"`
	//系统预留无法删除
	IsSysDefault bool `json:"isSysDefault"`
}
type UserRole struct {
	Id     uint `gorm:"primaryKey" json:"id"`
	UserId uint `json:"user_id"`
	RoleId uint `json:"role_id"`
}

type MenuRole struct {
	Id     uint `gorm:"primaryKey" json:"id"`
	RoleId uint `json:"role_id"`
	MenuId uint `json:"menu_id"`
}

type Resource struct {
	Id           uint   `gorm:"primaryKey" json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	MenuId       uint   `json:"menuId"`
	IsSysDefault bool   `json:"isSysDefault"`
}

type RoleResource struct {
	Id           uint   `gorm:"primaryKey" json:"id"`
	ResourcePath string `json:"resourcePath"`
	RoleId       uint   `json:"roleId"`
}

type History struct {
	Id           uint   `json:"id"`
	UserId       uint   `json:"userId"`
	Path         string `json:"path"`
	Method       string `json:"method"`
	FormData     string `json:"formData"`
	RequestType  string `json:"requestType"`
	RequestBody  string `json:"requestBody"`
	ResponseType string `json:"responseType"`
	ResponseBody string `json:"responseBody"`
}

func (h History) GetId() uint {
	return h.Id
}
