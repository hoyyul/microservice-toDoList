package jwts

import (
	"errors"
	"fmt"
	"micro-toDoList/global"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type CustomClaim struct {
	userId int64
	jwt.StandardClaims
}

func GenerateToken(userId int64) (string, error) {
	claim := CustomClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 3)), // expire time
			Issuer:    "xxxx",                                // issuer
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(getMySecret())
}

func ParseToken(tokenStr string) (*CustomClaim, error) {
	// 根据字符串取token
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) { return getMySecret(), nil })
	if err != nil {
		global.Logger.Error(fmt.Sprintf("token parse err: %s", err.Error()))
		return nil, err
	}

	// token还原claim， 返回claim
	if claim, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("invalid token")
}

func getMySecret() []byte {
	return []byte(global.Config.Server.Jwt)
}
