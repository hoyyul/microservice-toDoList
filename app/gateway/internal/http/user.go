package http

import (
	"micro-toDoList/app/gateway/internal/cache"
	"micro-toDoList/app/gateway/rpc"
	"micro-toDoList/global"
	"micro-toDoList/pkg/pb/user_pb"
	"micro-toDoList/pkg/resp"
	"micro-toDoList/pkg/util/jwts"

	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogin(ctx *gin.Context) {
	var req user_pb.UserRequest
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

	resp.SendWithOk(http.StatusOK, token, "Login successfully", ctx)
}

func UserRegister(ctx *gin.Context) {
	var req user_pb.UserRequest
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
	resp.SendWithOk(http.StatusOK, r, "Register a user successfully", ctx)
}

// no need to call grpc service and operate database.
func UserLogout(ctx *gin.Context) {
	_claim, _ := ctx.Get("claim")
	claim := _claim.(*jwts.CustomClaim)

	token := ctx.Request.Header.Get("token")

	// save expired token in cache
	cache.NewRedisService().SaveToken(*claim, token)

	resp.SendWithOk(http.StatusOK, map[string]interface{}{}, "User logout successfully", ctx)
}
