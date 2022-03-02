package history

import (
	"github.com/gin-gonic/gin"
	"go03/conf"
	"go03/db"
	"go03/es"
	"log"
	"net/http"
	"sync"
)

func List(ctx *gin.Context) {
	query := new(db.HistoryListQuery)
	err := ctx.BindJSON(query)
	if err != nil {
		ctx.JSON(http.StatusOK, conf.ServiceErrorResult("json转换失败"))
		return
	}
	ctx.JSON(http.StatusOK, db.HistoryDao.List(query))
	//ctx.JSON(http.StatusOK, es.HistoryEs.FindByPage(query.PageSize, query.Current))
}

func IndexAll(ctx *gin.Context) {
	maxId := db.HistoryDao.MaxId()
	var maxSize uint = 1000
	var w sync.WaitGroup
	for maxId >= 0 {
		hs := db.HistoryDao.FindByMaxIdPage(maxId, int(maxSize))
		if len(hs) <= 0 {
			break
		}
		w.Add(1)
		func() {
			defer func() {
				w.Done()
				err := recover()
				if err != nil {
					log.Println("index error:", err)
				}
			}()
			es.HistoryEs.IndexBatch(hs)
		}()
		if maxId <= maxSize {
			break
		}
		maxId = maxId - maxSize
	}
	w.Wait()
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}
