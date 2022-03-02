package auth

import "github.com/golang-jwt/jwt"

const (
	CookieName = "auth_token"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserClaims struct {
	jwt.StandardClaims
	UserId   uint   `json:"user_id"`
	Password string `json:"password"`
	Ip       string
	Agent    string
}

type LoginResult struct {
	UserId    uint   `json:"userId"`
	AuthToken string `json:"authToken"`
	Expires   int    `json:"expires"`
}

type UserVO struct {
	Userid uint `json:"userid"`
	//姓名
	Name string `json:"name"`
	//头像地址
	Avatar string `json:"avatar"`
	//邮箱
	Email string `json:"email"`
	//个性签名
	Signature string `json:"signature"`
	//头衔
	Title string `json:"title"`
	//地址
	Address string `json:"address"`
	//手机号
	Phone string `json:"phone"`
	//访问权限
	Access string `json:"access"`
}
