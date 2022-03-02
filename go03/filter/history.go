package filter

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go03/db"
	"go03/es"
	"io/ioutil"
	"log"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func safeSub(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[0:length]
}
func HistoryFilter(ctx *gin.Context) {
	reqBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBytes))
	delegateWriter := responseBodyWriter{
		ResponseWriter: ctx.Writer,
		body:           &bytes.Buffer{},
	}
	ctx.Writer = delegateWriter
	ctx.Next()
	val, _ := ctx.Get(gin.AuthUserKey)
	var uid uint
	if val != nil {
		u := val.(db.User)
		uid = u.Id
	}
	h := &db.History{
		Id:           0,
		UserId:       uid,
		Path:         ctx.Request.RequestURI,
		Method:       ctx.Request.Method,
		FormData:     "",
		RequestType:  ctx.Request.Header.Get("Content-Type"),
		RequestBody:  safeSub(string(reqBytes), 500),
		ResponseType: ctx.Writer.Header().Get("Content-Type"),
		ResponseBody: safeSub(delegateWriter.body.String(), 500),
	}
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("add history error:", err)
			}
		}()
		db.HistoryDao.Add(h)
		es.HistoryEs.IndexBatch([]*db.History{h})
	}()
}
