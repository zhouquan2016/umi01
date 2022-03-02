package db

import (
	"fmt"
	"go03/conf"
	//"go03/menu"
	"gorm.io/gorm"
	"strings"
)

type menuDao int

var MenuDao menuDao

func (m *menuDao) FindByUserId(userId uint) (menus []*Menu) {
	if userId <= 0 {
		return
	}
	u := UserDao.GetById(userId)
	if u == nil {
		return
	}

	role := RoleDao.FindByNo(u.RoleNo)
	if role == nil {
		return
	}
	menus = m.FindByRoleId(role.Id)
	menus = m.ToTree(menus)
	return
}

func (m menuDao) ToTree(menus []*Menu) []*Menu {
	var menuMap = map[uint][]*Menu{}
	var rootMenus = make([]*Menu, 0)
	for _, menu := range menus {
		if menu.ParentId == 0 {
			rootMenus = append(rootMenus, menu)
		} else {
			var children = menuMap[menu.ParentId]
			if children == nil {
				children = make([]*Menu, 0)
			}
			children = append(children, menu)
			menuMap[menu.ParentId] = children
		}
	}
	for _, parent := range rootMenus {
		var curArray = []*Menu{parent}
		for len(curArray) > 0 {
			var newArray = make([]*Menu, 0)
			for _, cur := range curArray {
				cur.Children = menuMap[cur.Id]
				if len(cur.Children) > 0 {
					newArray = append(newArray, cur.Children...)
				}
			}
			curArray = newArray
		}

	}
	return rootMenus
}
func (m *menuDao) FindByRoleId(id uint) []*Menu {
	var menus []*Menu
	query := db.Model(&Menu{}).Joins("join menu_role mr on mr.menu_id = menu.id and mr.role_id=?", id).Find(&menus)
	if query.Error != nil {
		panic(query.Error)
	}
	if query.RowsAffected <= 0 {
		return nil
	}
	return menus
}

type ListQuery struct {
	conf.BaseQuery
	ParentId     uint   `json:"parentId"`
	Path         string `json:"path"`
	IsLeaf       *bool  `json:"isLeaf,string"`
	IsSysDefault *bool  `json:"isSysDefault,string"`
}

type ListResult struct {
	Menus    []*Menu `json:"menus"`
	PageSize int     `json:"pageSize"`
	Current  int     `json:"current"`
	Total    int     `json:"total"`
}

func (menuDao *menuDao) FindPage(query *ListQuery) (page *conf.PageResult) {
	dbExe := db.Model(&Menu{})
	if query.ParentId > 0 {
		p := menuDao.FindById(query.ParentId)
		if p == nil {
			page = conf.NewPageResult(query.Current, 0, nil)
			return
		}
		dbExe.Where("dept like ? ", p.Dept+"%")
	}
	if query.Path != "" {
		dbExe.Where("path like ?", "%"+query.Path+"%")
	}
	if query.IsLeaf != nil {
		dbExe.Where("is_leaf = ? ", query.IsLeaf)
	}
	if query.IsSysDefault != nil {
		dbExe.Where("is_sys_default = ? ", query.IsSysDefault)
	}
	if query.SortFiled != "" {
		dbExe = dbExe.Order(query.SortFiled + " " + strings.ReplaceAll(query.SortOrder, "end", ""))
	}
	var total int64 = 0
	dbExe.Count(&total)
	var menus []*Menu
	if total > 0 {
		dbExe = dbExe.Offset(query.Offset()).Limit(query.PageSize).Find(&menus)
		if dbExe.Error != nil {
			panic(dbExe.Error)
		}
		if dbExe.RowsAffected <= 0 {
			menus = nil
		}
		var mids = make([]uint, 0)
		for _, menu := range menus {
			if menu.IsLeaf {
				mids = append(mids, menu.Id)
			}
		}
		rs := ResourceDao.FindByMenuIds(mids)
		for _, menu := range menus {
			menu.Resources = rs[menu.Id]
		}
	}
	page = conf.NewPageResult(query.Current, total, menus)
	return
}

func (dao *menuDao) FindAll() (menus []*Menu) {
	dbExe := db.Model(&Menu{}).Find(&menus)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		menus = nil
	}
	return
}

func (m *menuDao) Children(id uint) (menus []*Menu) {
	db.Model(&Menu{}).Where("parent_id", id).Find(&menus)
	return
}

func (m *menuDao) Add(menu Menu, createResources []*Resource) uint {

	menu.Id = 0
	var parentMenu = m.FindById(menu.ParentId)
	if menu.ParentId < 0 {
		panic(conf.ServiceErrorResult("父菜单id非法"))
	} else if menu.ParentId > 0 && parentMenu == nil {
		panic(conf.ServiceErrorResult("父菜单不存在"))
	} else if m.existsPath(menu.Path) {
		panic(conf.ServiceErrorResult("路径已存在"))
	}
	tx := db.Begin()
	defer func() {
		err := recover()
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			tx.Rollback()
			panic(err)

		}
	}()

	err := tx.Model(&Menu{}).Create(&menu).Error
	if err != nil {
		panic(tx.Error)
	}
	if parentMenu == nil {
		menu.Dept = fmt.Sprintf("%d:", menu.Id)
		menu.DeptName = "/" + menu.Name
	} else {
		menu.Dept = fmt.Sprintf("%s%d:", parentMenu.Dept, menu.Id)
		menu.DeptName = parentMenu.DeptName + "/" + menu.Name
	}

	err = tx.Model(&Menu{}).Where("id", menu.Id).Updates(map[string]interface{}{"dept_name": menu.DeptName, "dept": menu.Dept}).Error
	if err != nil {
		panic(tx.Error)
	}

	for _, resource := range createResources {
		resource.MenuId = menu.Id
		resource.Path = menu.Path + "/" + resource.Code
	}
	if len(createResources) > 0 {
		err = tx.Model(&Resource{}).Create(createResources).Error
		if err != nil {
			panic(err)
		}
	}

	return menu.Id
}

func (dao *menuDao) FindById(id uint) (menu *Menu) {
	dbExe := db.Model(&Menu{}).Where("id", id).Find(&menu)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		return nil
	}
	return
}

func (menuDao *menuDao) existsPath(path string) bool {
	dbExe := db.Model(&Menu{}).Where("path", path).First(&Menu{})
	if dbExe.Error != nil {
		if dbExe.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		return false
	}
	return true
}

func (menuDao *menuDao) HasChild(id uint) bool {
	dbExe := db.Model(&Menu{}).Where("parent_id", id).First(&Menu{})
	if dbExe.Error != nil {
		if dbExe.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		return false
	}
	return true
}

func (dao *menuDao) FindByIds(ids []uint) (menus []*Menu) {
	dbExe := db.Model(&Menu{}).Where("id in (?)", ids).Find(&menus)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	return
}

func (m *menuDao) Delete(ids []uint) {
	tx := db.Begin()
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
		if err != nil {
			panic(err)
		}
	}()
	err := tx.Model(&Menu{}).Where("id in (?)", ids).Delete(&Menu{}).Error
	if err != nil {
		panic(err)
	}
	err = tx.Model(&Resource{}).Where("menu_id in (?)", ids).Delete(&Resource{}).Error
	if err != nil {
		panic(err)
	}
}

func (m *menuDao) Update(id uint, params map[string]interface{}) {
	if len(params) <= 0 {
		panic(conf.ServiceErrorResult("更新字段为空"))
	}
	dbExe := db.Model(&Menu{}).Where("id", id).Updates(params)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
}

func (m *menuDao) GetByPath(path string) (menu *Menu) {
	dbExe := db.Model(&Menu{}).Where("path = ?", path).Find(&menu)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		menu = nil
	}
	return
}

func (m *menuDao) HasPath(roleId uint, path string) bool {
	var id int
	dbExe := db.Model(&Menu{}).Select("menu.id").Joins("inner join menu_role mr on mr.menu_id = menu.id and mr.role_id=?", roleId).Where("path=?", path).Scan(&id)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	return dbExe.RowsAffected > 0
}

func (dao *menuDao) Updates(m *Menu, createResources []*Resource, updateResources []*Resource) {
	tx := db.Begin()
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
		if err != nil {
			panic(err)
		}
	}()
	old := dao.FindById(m.Id)
	if old == nil {
		panic(conf.ServiceErrorResult("菜单不存在"))
	}
	deptName := ""
	parent := dao.FindById(m.ParentId)
	if parent != nil {
		deptName = parent.DeptName
	}
	deptName += "/" + m.Name
	err := tx.Model(&Menu{}).Where("id", m.Id).Updates(map[string]interface{}{
		"name":      m.Name,
		"path":      m.Path,
		"dept_name": deptName,
	}).Error
	if err != nil {
		panic(err)
	}
	if len(createResources) > 0 {
		err = tx.Model(&Resource{}).Create(createResources).Error
		if err != nil {
			panic(err)
		}
	}

	for _, resource := range updateResources {
		err = tx.Model(&Resource{}).Where("id", resource.Id).Updates(map[string]interface{}{
			"code": resource.Code,
			"name": resource.Name,
			"path": resource.Path,
		}).Error
		if err != nil {
			panic(err)
		}
	}

}
