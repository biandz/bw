package middwear

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Rsp struct{
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

//用户信息类，作为生成token的参数
type UserClaims struct {
	User    interface{}
	//jwt-go提供的标准claim
	jwt.StandardClaims
}

var (
	//自定义的token秘钥
	secret = []byte("16849841325189456f487")
	//token有效时间（2小时）
	effectTime = 2 * time.Hour
)

// 生成token
func GenerateToken(claims *UserClaims) (bool,string,string) {
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return false,"生成token失败",""
	}
	return true,"生成token成功",sign
}

//验证token
func JwtVerify(token string) error {
	if token == "" {
		return errors.New("token is not empty")
	}else{
		//验证token，并存储在请求中
		_, ok, msg := parseToken(token)
		if !ok{
			return errors.New(msg)
		}
	}
	return nil
}

// 解析Token
func parseToken(tokenString string) (*UserClaims, bool, string) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return &UserClaims{},false,"token is valid"
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return &UserClaims{},false,"parse token is failed"
	}
	return claims,true,""
}

// 更新token
func Refresh(tokenString string) string {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("token is valid")
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
	rst,msg, t := GenerateToken(claims)
	if !rst{
		return msg
	}
	return t
}

