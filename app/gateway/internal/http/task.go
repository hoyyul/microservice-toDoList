package http

import (
	"micro-toDoList/app/gateway/rpc"
	"micro-toDoList/global"
	"micro-toDoList/pkg/pb/task_pb"
	"micro-toDoList/pkg/resp"
	"micro-toDoList/pkg/util/jwts"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TaskCreate(ctx *gin.Context) {
	var req task_pb.TaskRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	// set login userId to task; user doesn't know self id
	_claim, ok := ctx.Get("claim")
	if !ok {
		resp.SendWithNotOk(http.StatusBadRequest, "User not found", ctx)
		return
	}
	claim := _claim.(*jwts.CustomClaim)
	req.UserId = claim.UserId

	// create task
	r, err := rpc.TaskCreate(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}

	resp.SendWithOk(http.StatusOK, r, "Created task succussfully!", ctx)
}

func TaskDelete(ctx *gin.Context) {
	var req task_pb.TaskRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	// set login userId to task
	_claim, ok := ctx.Get("claim")
	if !ok {
		resp.SendWithNotOk(http.StatusBadRequest, "User not found", ctx)
		return
	}
	claim := _claim.(*jwts.CustomClaim)
	req.UserId = claim.UserId

	// create task
	r, err := rpc.TaskDelete(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}

	resp.SendWithOk(http.StatusOK, r, "Deleted task succussfully!", ctx)
}

func TaskUpdate(ctx *gin.Context) {
	var req task_pb.TaskRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	// set login userId to task
	_claim, ok := ctx.Get("claim")
	if !ok {
		resp.SendWithNotOk(http.StatusBadRequest, "User not found", ctx)
		return
	}
	claim := _claim.(*jwts.CustomClaim)
	req.UserId = claim.UserId

	// create task
	r, err := rpc.TaskUpdate(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}

	resp.SendWithOk(http.StatusOK, r, "Updated task succussfully!", ctx)
}

func TaskShow(ctx *gin.Context) {
	var req task_pb.TaskRequest
	err := ctx.Bind(&req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusBadRequest, "Failed to bind request", ctx)
		return
	}

	// set login userId to task
	_claim, ok := ctx.Get("claim")
	if !ok {
		resp.SendWithNotOk(http.StatusBadRequest, "User not found", ctx)
		return
	}
	claim := _claim.(*jwts.CustomClaim)
	req.UserId = claim.UserId

	// create task
	r, err := rpc.TaskShow(ctx, &req)
	if err != nil {
		global.Logger.Panicln(err)
		resp.SendWithNotOk(http.StatusInternalServerError, "Failed to call User RPC service", ctx)
		return
	}
	resp.SendWithOk(http.StatusOK, r, "Get task details succussfully!", ctx)
}
