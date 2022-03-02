package menu

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go03/conf"
	"go03/db"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func List(ctx *gin.Context) {
	var sort = new(db.ListQuery)
	err := ctx.BindJSON(sort)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
	} else {
		pageResult := db.MenuDao.FindPage(sort)
		ctx.JSON(http.StatusOK, pageResult)
	}

}

func Children(ctx *gin.Context) {
	parentId, err := strconv.ParseInt(ctx.Query("parentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("参数非法"))
	}
	menus := db.MenuDao.Children(uint(parentId))
	treeNodes := make([]TreeNode, len(menus))
	for i, menu := range menus {
		treeNodes[i] = TreeNode{
			Id:     menu.Id,
			PId:    fmt.Sprintf("%d", menu.ParentId),
			Value:  menu.Id,
			Title:  menu.Name,
			IsLeaf: menu.IsLeaf,
		}
	}
	ctx.JSON(http.StatusOK, conf.SuccessResult(treeNodes))
}

func Add(ctx *gin.Context) {
	addQuery := &AddQuery{}
	err := ctx.BindJSON(addQuery)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	checkAdd(addQuery)
	createResources := make([]*db.Resource, 0)
	for _, resource := range addQuery.Resources {
		createResources = append(createResources, &db.Resource{
			Id:           0,
			Code:         resource.Code,
			Name:         resource.Name,
			IsSysDefault: false,
		})
	}
	id := db.MenuDao.Add(db.Menu{
		Id:       0,
		Name:     *addQuery.Name,
		Children: nil,
		IsLeaf:   *addQuery.IsLeaf,
		Path:     *addQuery.Path,
		ParentId: *addQuery.ParentId,
	}, createResources)
	ctx.JSON(http.StatusOK, conf.SuccessResult(id))
}

func chekPath(path *string) bool {
	if path == nil || *path == "" {
		return false
	}
	if ok, err := regexp.Match("^(/[\\w_-]+)+$", []byte(*path)); !ok || err != nil {
		if err != nil {
			panic(err)
		}
		return false
	}
	if len([]rune(*path)) > 20 {
		return false
	}
	return true
}
func checkAdd(query *AddQuery) {
	errMsg := ""
	if query.ParentId == nil {
		errMsg = "未选择父菜单"
	} else if *query.ParentId < 0 {
		errMsg = "父菜单非法"
	} else if query.Name == nil || *query.Name == "" {
		errMsg = "菜单名称为空"
	} else if len([]rune(*query.Name)) > 10 {
		errMsg = "菜单名称长度不能大于10"
	} else if !chekPath(query.Path) {
		errMsg = "路径非法"
	} else if query.IsLeaf == nil {
		errMsg = "是否叶子节点为空"
	}
	if errMsg != "" {
		panic(conf.ServiceErrorResult(errMsg))
	}

}

func Delete(ctx *gin.Context) {
	var ids []uint
	err := ctx.BindJSON(&ids)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
	}
	menus := db.MenuDao.FindByIds(ids)
	deleteIds := make([]uint, 0)
	errList := make([]string, 0)
	for _, menu := range menus {
		if menu.IsSysDefault {
			errList = append(errList, menu.Name+"是系统预留菜单无法删除")
		} else if !menu.IsLeaf && db.MenuDao.HasChild(menu.Id) {
			errList = append(errList, menu.Name+"下还有子菜单，请先删除")
		} else {
			deleteIds = append(deleteIds, menu.Id)
		}

	}
	if len(errList) > 0 {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult(strings.Join(errList, ",")))
	} else {
		db.MenuDao.Delete(deleteIds)
		ctx.JSON(http.StatusOK, conf.SuccessResult(true))
	}

}

func Edit(ctx *gin.Context) {
	query := new(UpdateQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	checkEdit(query)
	m := db.MenuDao.FindById(query.Id)
	if m == nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("菜单不存在"))
		return
	}
	if m.IsSysDefault {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("菜单是系统预留无法修改"))
		return
	}
	m.Name = query.Name
	m.Path = query.Path
	updateResources := make([]*db.Resource, 0)
	createResources := make([]*db.Resource, 0)
	if m.IsLeaf {
		for _, resource := range query.Resources {
			r := &db.Resource{
				Id:           resource.Id,
				Code:         resource.Code,
				Name:         resource.Name,
				Path:         query.Path + "/" + resource.Code,
				MenuId:       query.Id,
				IsSysDefault: false,
			}
			if resource.Id <= 0 {
				createResources = append(createResources, r)
			} else {
				updateResources = append(updateResources, r)
			}
		}
	}
	db.MenuDao.Updates(m, createResources, updateResources)
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}

func checkEdit(query *UpdateQuery) {
	var errString = ""
	if query.Id <= 0 {
		errString = "id非法"
	} else if query.Name == "" {
		errString = "名称不能为空"
	} else if !chekPath(&query.Path) {
		errString = "路径非法"
	} else if db.MenuDao.FindById(query.Id) == nil {
		errString = "菜单不存在"
	}
	if errString != "" {
		panic(conf.ServiceErrorResult(errString))
	}

}

func GetById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("id非法"))
		return
	}
	var menu *db.Menu
	if id > 0 {
		uid := uint(id)
		menu = db.MenuDao.FindById(uid)
		rm := db.ResourceDao.FindByMenuIds([]uint{uid})
		if rm != nil {
			menu.Resources = rm[uid]
		}
	}
	ctx.JSON(http.StatusOK, conf.SuccessResult(menu))
}

func ExistsByPath(ctx *gin.Context) {
	var menu *db.Menu
	if path, ok := ctx.GetQuery("path"); path != "" && ok {
		menu = db.MenuDao.GetByPath(path)
	}
	ctx.JSON(http.StatusOK, conf.SuccessResult(menu != nil))
}

func AddResource(ctx *gin.Context) {
	addQuery := new(AddResourceQuery)
	err := ctx.BindJSON(addQuery)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}

	menuEntity := db.MenuDao.FindById(addQuery.MenuId)
	if menuEntity == nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("菜单不存在"))
		return
	}
	if menuEntity.IsSysDefault {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("菜单是系统预留无法删除修改资源"))
		return
	}
	var codeMap = make(map[string]*db.Resource)
	for _, code := range addQuery.Codes {
		if code != "" {
			codeMap[code] = &db.Resource{
				Id:           0,
				Code:         code,
				Path:         menuEntity.Path + "/" + code,
				MenuId:       menuEntity.Id,
				IsSysDefault: false,
			}
		}
	}
	var rs = make([]*db.Resource, 0)
	for _, resource := range codeMap {
		rs = append(rs, resource)
	}
	db.ResourceDao.Add(rs)
}

func Tree(ctx *gin.Context) {
	menus := db.MenuDao.FindAll()
	mids := make([]uint, 0)
	for _, menu := range menus {
		if menu.IsLeaf {
			mids = append(mids, menu.Id)
		}
	}
	rm := db.ResourceDao.FindByMenuIds(mids)
	if len(rm) > 0 {
		for _, menu := range menus {
			menu.Resources = rm[menu.Id]
		}
	}

	menus = db.MenuDao.ToTree(menus)
	ctx.JSON(http.StatusOK, conf.SuccessResult(menus))
}
