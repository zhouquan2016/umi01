package db

import (
	"go03/conf"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type roleDao int

var RoleDao roleDao

func (roleDao *roleDao) FindByNo(no string) (r *Role) {
	query := db.Model(&Role{}).Where("no", no).Find(&r)
	if err := query.Error; err != nil {
		panic(err)
	}
	if query.RowsAffected <= 0 {
		r = nil
	}
	return
}

type RoleListQuery struct {
	No           string `json:"no"`
	Name         string `json:"name"`
	IsSysDefault *bool  `json:"isSysDefault,string"`
	conf.BaseQuery
}

func (roleDao *roleDao) FindByPage(query *RoleListQuery) (roles []*Role, total int64) {
	dbQuery := db.Model(&Role{})
	if query.No != "" {
		dbQuery.Where("no like ?", "%"+query.No+"%")
	}
	if query.Name != "" {
		dbQuery.Where("name like ?", "%"+query.Name+"%")
	}
	if query.IsSysDefault != nil {
		dbQuery.Where("is_sys_default", query.IsSysDefault)
	}
	dbQuery.Offset(query.Offset()).Limit(query.PageSize)
	err := dbQuery.Count(&total).Error
	if err != nil {
		panic(err)
	}
	err = dbQuery.Find(&roles).Error
	if err != nil {
		panic(err)
	}
	return
}

func (dao *roleDao) Add(r *Role, resourceIds []uint) {
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
	err := tx.Model(&Role{}).Create(r).Error
	if err != nil {
		panic(err)
	}
	if len(resourceIds) > 0 {
		createMenuResource(tx, resourceIds, r.Id)
	}
}

func (dao *roleDao) DeleteById(id uint) {
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
	dbExe := tx.Exec("delete from role where id = ? and not exists(select 1 from user where role_id = role.id)", id)
	err := dbExe.Error
	if err != nil {
		panic(err)
	}
	if dbExe.RowsAffected <= 0 {
		panic(conf.ServiceErrorResult("删除角色失败"))
	}
	err = tx.Model(&RoleResource{}).Where("role_id", id).Delete(&RoleResource{}).Error
	if err != nil {
		panic(err)
	}
}

func (dao *roleDao) Update(r *Role, resourceIds []uint) {
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
	err := tx.Model(&Role{}).Where("id", r.Id).Updates(map[string]interface{}{"name": r.Name}).Error
	if err != nil {
		panic(err)
	}
	err = tx.Model(&RoleResource{}).Where("role_id", r.Id).Delete(&RoleResource{}).Error
	if err != nil {
		panic(err)
	}
	if len(resourceIds) <= 0 {
		return
	}
	err = tx.Model(&MenuRole{}).Where("role_id", r.Id).Delete(&MenuRole{}).Error
	if err != nil {
		panic(err)
	}
	createMenuResource(tx, resourceIds, r.Id)

}

func createMenuResource(tx *gorm.DB, resourceIds []uint, roleId uint) {
	rs := ResourceDao.FindByIds(resourceIds)
	rr := make([]*RoleResource, 0)
	mids := make([]uint, 0)
	for _, resource := range rs {
		rr = append(rr, &RoleResource{
			Id:           0,
			ResourcePath: resource.Path,
			RoleId:       roleId,
		})
		mids = append(mids, resource.MenuId)
	}
	err := tx.Model(&RoleResource{}).Create(rr).Error
	if err != nil {
		panic(err)
	}
	menus := MenuDao.FindByIds(mids)
	menuMap := make(map[uint]interface{}, 0)
	for _, menu := range menus {
		for _, sid := range strings.Split(menu.Dept, ":") {
			if sid == "" {
				continue
			}
			mid, err := strconv.Atoi(sid)
			if err != nil {
				panic(err)
			}
			menuMap[uint(mid)] = ""
		}
	}
	mrs := make([]*MenuRole, 0)
	for mid, _ := range menuMap {
		mrs = append(mrs, &MenuRole{
			Id:     0,
			RoleId: roleId,
			MenuId: mid,
		})
	}
	if len(mrs) > 0 {
		err := tx.Model(&MenuRole{}).Create(mrs).Error
		if err != nil {
			panic(err)
		}
	}
}

func (dao *roleDao) FindById(id uint) (r *Role) {
	if id <= 0 {
		return
	}
	dbQuery := db.Model(&Role{}).Where("id", id).Find(&r)
	if dbQuery.Error != nil {
		panic(dbQuery.Error)
	}
	if dbQuery.RowsAffected <= 0 {
		r = nil
	}
	return
}

func (dao *roleDao) FindAll() (roles []*Role) {
	err := db.Model(&Role{}).Find(&roles).Error
	if err != nil {
		panic(err)
	}
	return
}

func (dao *roleDao) FindByIds(ids []uint) (roleMap map[uint]*Role) {
	var roles []*Role
	err := db.Model(&Role{}).Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		panic(err)
	}
	roleMap = make(map[uint]*Role, 0)
	for _, r := range roles {
		roleMap[r.Id] = r
	}
	return
}
