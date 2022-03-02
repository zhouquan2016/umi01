package user

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go03/conf"
	"go03/db"
	"net/http"
	"strconv"
)

func FindMenus(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists || u == nil {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "未登录"))
	}
	var user = u.(db.User)
	ctx.JSON(http.StatusOK, conf.SuccessResult(db.MenuDao.FindByUserId(user.Id)))

}

type PageVo struct {
	Id uint `json:"id"`
	//姓名
	Name string `json:"name"`
	//头像地址
	Avatar string `json:"avatar"`
	//邮箱
	Email string `json:"email"`
	//个性签名
	Signature string `json:"signature"`
	//头衔
	Title string `json:"title"`
	//地址
	Address string `json:"address"`
	//手机号
	Phone string `json:"phone"`
	//角色id
	RoleId uint `json:"roleId"`
	//角色名称
	RoleName     string `json:"roleName"`
	IsSysDefault bool   `json:"isSysDefault"`
}

func List(ctx *gin.Context) {
	query := new(db.UserListQuery)
	if err := ctx.BindJSON(query); err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	toal, users := db.UserDao.List(query)
	var pageVos = make([]*PageVo, len(users))
	roleIds := make([]uint, 0)
	for _, u := range users {
		roleIds = append(roleIds, u.RoleId)
	}
	roleMap := db.RoleDao.FindByIds(roleIds)
	for i, user := range users {
		vo := &PageVo{
			Id:           user.Id,
			Name:         user.Name,
			Avatar:       user.Avatar,
			Email:        user.Email,
			Signature:    user.Signature,
			Title:        user.Title,
			Address:      user.Address,
			Phone:        user.Phone,
			IsSysDefault: user.IsSysDefault,
			RoleId:       user.RoleId,
		}
		if roleMap[user.RoleId] != nil {
			vo.RoleName = roleMap[user.RoleId].Name
		}
		pageVos[i] = vo
	}
	ctx.JSON(http.StatusOK, conf.NewPageResult(query.Current, toal, pageVos))
}

type AddQuery struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	RoleId    uint   `json:"roleId"`
}

func makeUser(query *AddQuery) (u *db.User) {
	r := db.RoleDao.FindById(query.RoleId)
	if r == nil {
		panic(conf.ServiceErrorResult("角色未找到"))
	}
	if db.UserDao.GetByEmail(query.Email) != nil {
		panic(conf.ServiceErrorResult("邮箱已注册"))
	}
	if db.UserDao.GetByPhone(query.Phone) != nil {
		panic(conf.ServiceErrorResult("手机号已注册"))
	}
	pwd := md5.Sum([]byte(query.Password))
	u = &db.User{
		Id:           0,
		Name:         query.Name,
		Password:     hex.EncodeToString(pwd[:]),
		Avatar:       "",
		Email:        query.Email,
		Signature:    query.Signature,
		Title:        query.Title,
		Address:      query.Address,
		Phone:        query.Phone,
		RoleNo:       r.No,
		RoleId:       r.Id,
		IsSysDefault: false,
	}
	return
}
func Save(ctx *gin.Context) {
	query := new(AddQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	u := makeUser(query)
	db.UserDao.Add(u)
}

type EditQuery struct {
	AddQuery
	Id uint `json:"id"`
}

func Edit(ctx *gin.Context) {
	query := new(EditQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	u := db.UserDao.GetById(query.Id)
	if db.UserDao.GetById(query.Id) == nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("用户不存在"))
		return
	}
	if u.IsSysDefault {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("系统预留无法修改"))
		return
	}
	r := db.RoleDao.FindById(query.RoleId)
	if r == nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("角色不存在"))
		return
	}
	u.Name = query.Name
	u.Signature = query.Signature
	u.Title = query.Title
	u.Address = query.Address
	u.RoleId = r.Id
	u.RoleNo = r.No
	db.UserDao.Update(u)
}

func parseId(ctx *gin.Context) (uint, bool) {
	idStr, _ := ctx.GetQuery("id")
	if idStr == "" {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("id为空"))
		return 0, false
	}
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("id非法"))
		return 0, false
	}
	return uint(id), true
}
func GetById(ctx *gin.Context) {
	id, ok := parseId(ctx)
	if !ok {
		return
	}
	u := db.UserDao.GetById(uint(id))
	ctx.JSON(http.StatusOK, conf.SuccessResult(&PageVo{
		Id:        u.Id,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Signature: u.Signature,
		Title:     u.Title,
		Address:   u.Address,
		Phone:     u.Phone,
		RoleId:    u.RoleId,
	}))
}

func DeleteById(ctx *gin.Context) {
	id, ok := parseId(ctx)
	if !ok {
		return
	}
	db.UserDao.DeleteById(id)
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}
