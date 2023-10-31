package http

import (
	"go-micro-toDoList/app/gateway/rpc"
	"go-micro-toDoList/global"
	"go-micro-toDoList/pkg/pb"
	"go-micro-toDoList/pkg/resp"
	"go-micro-toDoList/pkg/util/jwts"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	var req pb.UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	r, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}

	// generate jwt
	token, err := jwts.GenerateToken(r.UserId)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to generate jwt token", ctx)
		return
	}

	resp.SendWithOk(http.StatusOK, token, ctx)
}

func UserRegister(ctx *gin.Context) {
	var req pb.UserRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	r, err := rpc.UserRegister(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}
	resp.SendWithOk(http.StatusOK, r, ctx)
}
