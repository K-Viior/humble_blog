package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 自定义声明结构体并内嵌jwt.StandardClaims
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var secret string

func initSecret() {
	secret = "hello"
}

func generateToken(userId int64, expire time.Duration) (string, error) {
	if secret == "" {
		initSecret()
	}
	mySigningKey := []byte(secret)
	// Create the claims
	claims := CustomClaims{
		userId,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),             // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),             // 生效时间
			Issuer:    "novo",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	//fmt.Printf("%v %v", tokenString, err)
	return tokenString, err
}

func AccessToken(userId int64) (string, error) {
	return generateToken(userId, 2*time.Hour)
}

func RefreshToken(userId int64) (string, error) {
	return generateToken(userId, 7*time.Hour)
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	if secret == "" {
		initSecret()
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("解析token失败")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.UserID, claims.RegisteredClaims.Issuer)
		return claims, nil
	}
	return nil, errors.New("token不合法")
}

func main() {
	token, err := AccessToken(23)
	if err != nil {

	}
	fmt.Println(token)
	claims, err := VerifyToken(token)
	fmt.Printf("%v", claims)
}
