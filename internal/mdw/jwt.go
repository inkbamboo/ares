// Package mdw 提供中间件功能。
package mdw

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// JwtHandler 是 JWT 认证中间件的占位实现（尚未完成）。
func JwtHandler() {
	a := jwt.Token{}
	fmt.Println(a)
	//jwt.New()
}
