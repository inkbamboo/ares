package mdw

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

/**
 * jwt中间件
 */
func JwtHandler() {
	a := jwt.Token{}
	fmt.Println(a)
	//jwt.New()
}
