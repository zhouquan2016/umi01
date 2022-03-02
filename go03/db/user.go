package db

import (
	"go03/conf"
	"strings"
)

type userDao int

var UserDao userDao

type UserListQuery struct {
	conf.BaseQuery
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	RoleId uint   `json:"roleId"`
}

func (userDao *userDao) GetByEmail(email string) (u *User) {
	dbQuery := db.Model(&User{}).Where("email", email).Find(&u)
	if dbQuery.Error != nil {
		panic(dbQuery.Error)
	}
	if dbQuery.RowsAffected <= 0 {
		u = nil
	}
	return
}

func (userDao *userDao) GetById(id uint) (u *User) {
	if id <= 0 {
		return nil
	}
	dbQuery := db.Find(&u, id)
	if dbQuery.Error != nil {
		panic(dbQuery.Error)
	}
	if dbQuery.RowsAffected <= 0 {
		u = nil
	}
	return
}

func (userDao *userDao) GetByPhone(phone string) (u *User) {
	dbQuery := db.Model(&User{}).Where("phone", phone).Find(&u)
	if dbQuery.Error != nil {
		panic(dbQuery.Error)
	}
	if dbQuery.RowsAffected <= 0 {
		u = nil
	}
	return
}

func (userDao *userDao) List(query *UserListQuery) (total int64, users []*User) {
	dbExe := db.Model(&User{})
	if query.Name != "" {
		dbExe = dbExe.Where("name like ?", "%"+query.Name+"%")
	}
	if query.Email != "" {
		dbExe = dbExe.Where("email like ?", "%"+query.Email+"%")
	}
	if query.Phone != "" {
		dbExe = dbExe.Where("phone like ?", "%"+query.Phone+"%")
	}
	if query.RoleId > 0 {
		dbExe = dbExe.Where("role_id", query.RoleId)
	}
	if query.SortFiled != "" {
		dbExe = dbExe.Order(query.SortFiled + " " + strings.ReplaceAll(query.SortOrder, "end", ""))
	}
	dbExe = dbExe.Count(&total)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if total > 0 {
		dbExe = dbExe.Offset(query.Offset()).Limit(query.PageSize).Find(&users)
		if dbExe.Error != nil {
			panic(dbExe.Error)
		}
		if dbExe.RowsAffected <= 0 {
			users = nil
		}
	}
	return
}

func (dao *userDao) Add(u *User) {
	if u == nil {
		return
	}
	err := db.Model(&User{}).Create(u).Error
	if err != nil {
		panic(err)
	}
}

func (dao *userDao) Update(u *User) {
	if u == nil {
		return
	}
	err := db.Model(&User{}).Where("id", u.Id).UpdateColumns(u).Error
	if err != nil {
		panic(err)
	}
}

func (dao *userDao) DeleteById(id uint) {
	u := dao.GetById(id)
	if u == nil {
		panic(conf.ServiceErrorResult("用户不存在"))
	}
	if u.IsSysDefault {
		panic(conf.ServiceErrorResult("系统预留无法删除"))
	}
	err := db.Model(&User{}).Where("id", id).Delete(&User{}).Error
	if err != nil {
		panic(err)
	}
}

func (dao *userDao) ExistsByRoleId(id uint) bool {
	var exists bool
	err := db.Raw("select exists(select 1 from user where role_id = ?)", id).Scan(&exists).Error
	if err != nil {
		panic(err)
	}
	return exists

}
