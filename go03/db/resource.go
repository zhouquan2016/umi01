package db

type resourceDao uint

var ResourceDao resourceDao

func (dao *resourceDao) Add(rs []*Resource) {
	dbExe := db.Model(&Resource{}).Create(rs)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
}

func (dao *resourceDao) FindByPath(path string) (r *Resource) {
	dbExe := db.Model(&Resource{}).Where("path", path).Find(&r)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		r = nil
	}
	return
}

func (dao *resourceDao) FindByMenuIds(menuIds []uint) map[uint][]*Resource {
	var rs []*Resource
	dbExe := db.Model(&Resource{}).Where("menu_id in (?)", menuIds).Find(&rs)
	if dbExe.Error != nil {
		panic(dbExe.Error)
	}
	if dbExe.RowsAffected <= 0 {
		return nil
	}
	var m = make(map[uint][]*Resource)
	for _, r := range rs {
		s := m[r.MenuId]
		if s == nil {
			s = make([]*Resource, 0)
		}
		s = append(s, r)
		m[r.MenuId] = s
	}
	return m
}

func (dao *resourceDao) FindByIds(ids []uint) (rs []*Resource) {
	err := db.Model(&Resource{}).Where("id in ?", ids).Find(&rs).Error
	if err != nil {
		panic(err)
	}
	return
}

func (dao *resourceDao) FindByRoleId(roleId uint) (rs []*Resource) {
	err := db.Model(&Resource{}).Joins("inner join role_resource r on r.resource_path = resource.path and r.role_id = ?", roleId).Find(&rs).Error
	if err != nil {
		panic(err)
	}
	return
}

func (dao *resourceDao) ExistsRolePath(roleId uint, path string) bool {
	dbQuery := db.Model(&RoleResource{}).Where("role_id", roleId).Where("resource_path", path).Find(&RoleResource{})
	if dbQuery.Error != nil {
		panic(dbQuery.Error)
	}
	return dbQuery.RowsAffected > 0
}
