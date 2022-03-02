package conf

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime/debug"
)

const (
	SuccessStatus      = 1
	ServiceErrorStatus = 2
	SystemErrorStatus  = 3
	LoginErrorStatus   = 4
)

type ApiResult struct {
	Success      bool        `json:"success"`
	Status       int8        `json:"status"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

//成功返回
func SuccessResult(val interface{}) *ApiResult {
	return &ApiResult{
		Success:      true,
		Status:       SuccessStatus,
		ErrorMessage: "成功",
		Data:         val,
	}
}

//失败返回
func ErrorResult(status int8, message string) *ApiResult {
	return &ApiResult{
		Success:      false,
		Status:       status,
		ErrorMessage: message,
		Data:         nil,
	}
}

//业务异常
func ServiceErrorResult(message string) *ApiResult {
	return ErrorResult(ServiceErrorStatus, message)
}

//系统异常
func SysErrorResult() *ApiResult {
	return ErrorResult(ServiceErrorStatus, "系统异常")
}

func ErrorFilter(context *gin.Context) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		r, ok := err.(*ApiResult)
		if !ok {
			catchErr, ok := err.(error)
			if ok && catchErr != nil {
				log.Println("Catch Error:", catchErr.Error())
			} else {
				log.Println("Catch Error:", catchErr)
			}

			debug.PrintStack()
			r = SysErrorResult()
		} else {
			log.Println("service error:", err)
		}

		context.AbortWithStatusJSON(http.StatusOK, r)
	}()
	context.Next()
}

func Assert(expect bool, message string) {
	if !expect {
		panic(ServiceErrorResult(message))
	}
}
