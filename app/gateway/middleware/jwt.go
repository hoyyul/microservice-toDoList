package middleware

import (
	"micro-toDoList/pkg/resp"
	"micro-toDoList/pkg/util/jwts"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")

		// check if empty
		if token == "" {
			resp.SendWithNotOk(http.StatusOK, "Request is missing a token", ctx)
			ctx.Abort()
			return
		}

		// check if token valid; return error if incorrect or expired
		claim, err := jwts.ParseToken(token)
		if err != nil {
			resp.SendWithNotOk(http.StatusOK, "No information found for the given token", ctx)
			ctx.Abort()
			return
		}

		// check if in redis
		/*if redis_service.CheckLogout(token) {
			res.FailWithMessage("token expired", c)
			ctx.Abort()
			return
		}*/

		// save user infor in gin
		ctx.Set("claim", claim)
		ctx.Next()
	}
}
