package resp

import (
	"micro-toDoList/pkg/errmsg"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func SendWithOk(status int, data interface{}, msg string, ctx *gin.Context) {
	ctx.JSON(status, Response{
		Code: errmsg.SUCCESS,
		Data: data,
		Msg:  "",
	})
}

func SendWithNotOk(status int, msg string, ctx *gin.Context) {
	ctx.JSON(status, Response{
		Code: errmsg.FAILURE,
		Data: map[string]interface{}{}, // json传空值是这个格式
		Msg:  msg,
	})
}
