package db

import "go03/conf"

type historyDao int

var HistoryDao historyDao

func (dao historyDao) Add(h *History) {
	if h != nil {
		err := db.Model(&History{}).Create(h).Error
		if err != nil {
			panic(err)
		}
	}
}

type HistoryListQuery struct {
	conf.BaseQuery
	History
}

func (dao historyDao) List(query *HistoryListQuery) *conf.PageResult {
	dbQuery := db.Model(&History{})
	var total int64 = 0
	err := dbQuery.Count(&total).Error
	if err != nil {
		panic(err)
	}
	var hs []*History
	if total > 0 {
		if query.UserId != 0 {
			dbQuery = dbQuery.Where("user_id", query.UserId)
		}
		if query.Path != "" {
			dbQuery = dbQuery.Where("path like ?", "%"+query.Path+"%")
		}
		if query.Method != "" {
			dbQuery = dbQuery.Where("method", query.Method)
		}
		if query.FormData != "" {
			dbQuery = dbQuery.Where("form_data like ?", "%"+query.FormData+"%")
		}
		if query.RequestType != "" {
			dbQuery = dbQuery.Where("request_type like ?", "%"+query.RequestType+"%")
		}
		if query.RequestBody != "" {
			dbQuery = dbQuery.Where("request_body like ?", "%"+query.RequestBody+"%")
		}
		if query.ResponseType != "" {
			dbQuery = dbQuery.Where("response_type like ?", "%"+query.ResponseType+"%")
		}
		if query.RequestBody != "" {
			dbQuery = dbQuery.Where("request_body like ?", "%"+query.RequestBody+"%")
		}
		err := dbQuery.Offset(query.Offset()).Limit(query.PageSize).Order("id desc").Find(&hs).Error
		if err != nil {
			panic(err)
		}
	}
	return conf.NewPageResult(query.Current, total, hs)
}

func (dao historyDao) MaxId() (id uint) {
	err := db.Model(&History{}).Select("id").Order("id desc").Limit(1).Scan(&id).Error
	if err != nil {
		panic(err)
	}
	return
}

func (dao historyDao) FindByMaxIdPage(maxId uint, maxSize int) (hs []*History) {
	err := db.Model(&History{}).Where("id <= ?", maxId).Order("id desc").Limit(maxSize).Find(&hs).Error
	if err != nil {
		panic(err)
	}
	return
}
