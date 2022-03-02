package role

import (
	"github.com/gin-gonic/gin"
	"go03/conf"
	"go03/db"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	query := new(db.RoleListQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换异常"))
		return
	}
	roles, total := db.RoleDao.FindByPage(query)
	page := conf.NewPageResult(query.Current, total, roles)
	ctx.JSON(http.StatusOK, page)
}

func Exists(ctx *gin.Context) {
	no, _ := ctx.GetQuery("no")
	if no == "" {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("参数为空"))
		return
	}
	r := db.RoleDao.FindByNo(no)

	ctx.JSON(http.StatusOK, conf.SuccessResult(r != nil))
}

type AddQuery struct {
	No          string `json:"no"`
	Name        string `json:"name"`
	ResourceIds []uint `json:"resourceIds"`
}

func Add(ctx *gin.Context) {
	query := new(AddQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换异常"))
		return
	}
	checkAdd(query)
	r := &db.Role{
		Id:           0,
		No:           query.No,
		Name:         query.Name,
		IsSysDefault: false,
	}
	db.RoleDao.Add(r, query.ResourceIds)
}

func checkAdd(query *AddQuery) {
	if len(query.No) < 2 || len(query.No) > 10 {
		panic(conf.ServiceErrorResult("角色编码长度应该在2~5之间"))
	}
	if len(query.Name) < 2 || len(query.Name) > 20 {
		panic(conf.ServiceErrorResult("角色名称长度应该在2~10之间"))
	}
	r := db.RoleDao.FindByNo(query.No)
	if r != nil {
		panic(conf.ServiceErrorResult("角色编码已存在"))
	}
}

func Delete(ctx *gin.Context) {
	id, ok := getQueryId(ctx)
	if !ok {
		return
	}
	r := db.RoleDao.FindById(uint(id))
	if r == nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("角色不存在"))
		return
	}
	if db.UserDao.ExistsByRoleId(r.Id) {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("当前角色下已有关联用户，无法删除"))
		return
	}
	db.RoleDao.DeleteById(uint(id))
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}

type UpdateQuery struct {
	Name        string `json:"name"`
	Id          uint   `json:"id"`
	ResourceIds []uint `json:"resourceIds"`
}

func Update(ctx *gin.Context) {
	query := new(UpdateQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		panic(err)
	}
	db.RoleDao.Update(&db.Role{
		Id:   query.Id,
		Name: query.Name,
	}, query.ResourceIds)
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}

func getQueryId(ctx *gin.Context) (uint, bool) {
	idStr, _ := ctx.GetQuery("id")
	if idStr == "" {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("id为空"))
		return 0, false
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("id非法"))
		return 0, false
	}
	return uint(id), true
}

type BaseRoleVo struct {
	db.Role
	ResourceIds []uint `json:"resourceIds"`
}

func GetById(ctx *gin.Context) {
	id, ok := getQueryId(ctx)
	if !ok {
		return
	}
	r := db.RoleDao.FindById(id)
	var vo *BaseRoleVo
	if r != nil {
		vo = new(BaseRoleVo)
		vo.Role = *r
		resourceIds := make([]uint, 0)
		rr := db.ResourceDao.FindByRoleId(r.Id)
		for _, resource := range rr {
			resourceIds = append(resourceIds, resource.Id)
		}
		vo.ResourceIds = resourceIds
	}
	ctx.JSON(http.StatusOK, conf.SuccessResult(vo))
}

func GetAll(ctx *gin.Context) {
	roles := db.RoleDao.FindAll()
	ctx.JSON(http.StatusOK, conf.SuccessResult(roles))
}
