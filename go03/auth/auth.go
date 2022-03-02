package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go03/conf"
	"go03/db"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const LoginPath = "/login/account"
const LogoutPath = "/login/outLogin"

//生成登录token
func genJwtToken(u *db.User, r *http.Request) string {
	mySigningKey := []byte(conf.GetConfig().JwtSecret)

	requestClaims := pareRequestClaims(r)
	log.Println("jwt ip:", requestClaims.Ip)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		},
		UserId:   u.Id,
		Password: u.Password,
		Ip:       requestClaims.Ip,
		Agent:    requestClaims.Agent,
	})
	tokenString, err := token.SignedString(mySigningKey)
	conf.Assert(err == nil, "生成token失败")
	return tokenString
}

func pareRequestClaims(r *http.Request) *UserClaims {
	host := r.Header.Get("X-Real-Ip")
	if host == "" {
		host, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return &UserClaims{Ip: host, Agent: r.Header.Get("User-Agent")}
}

//校验登录的token
func validateJwtToken(tokenString string, r *http.Request) *db.User {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetConfig().JwtSecret), nil
	})
	if err != nil {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "token非法"))
	}
	claims := token.Claims.(*UserClaims)
	u := db.UserDao.GetById(claims.UserId)
	if u == nil {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "用户不存在"))
	}
	if u.Password != claims.Password {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "密码不正确"))
	}
	requestClaims := pareRequestClaims(r)
	if !net.ParseIP(claims.Ip).Equal(net.ParseIP(requestClaims.Ip)) {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "网络环境已发生变化"))
	}
	if claims.Agent != requestClaims.Agent {
		panic(conf.ErrorResult(conf.LoginErrorStatus, "客户端已发生变化"))
	}
	return u
}

//获取登录的cookie
func getCookie(ctx *gin.Context) *http.Cookie {
	cookie, err := ctx.Request.Cookie(CookieName)
	if err != nil {
		return nil
	}
	return cookie
}

//请求拦截需要登录
func Filter(context *gin.Context) {
	if context.Request.RequestURI == LoginPath {
		return
	}
	cookie := getCookie(context)

	if cookie == nil || cookie.Value == "" {
		context.AbortWithStatusJSON(http.StatusOK, conf.ErrorResult(conf.LoginErrorStatus, "未登录"))
		return
	}
	u := validateJwtToken(cookie.Value, context.Request)
	if u != nil {
		context.Set("user", *u)
		context.Next()
	}

}

func getBrowerDomain(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin != "" {
		origin = origin[strings.Index(origin, "://")+3:]
		index := strings.Index(origin, "://")
		if index > 0 {
			origin = origin[index+3:]
		}
		index = strings.LastIndex(origin, "/")
		if index > 0 {
			origin = origin[:index]
		}
		return origin
	}
	domain, _, _ := net.SplitHostPort(r.Host)
	return domain
}

//登录请求
func Login(context *gin.Context) {
	loginUser := new(LoginUser)
	_ = context.BindJSON(loginUser)
	conf.Assert(loginUser.Username != "", "用户名不能为空")
	conf.Assert(loginUser.Password != "", "密码不能为空")
	u := findUser(loginUser.Username)
	conf.Assert(u != nil, "用户不存在")

	context.SetCookie(CookieName, genJwtToken(u, context.Request), 60*60*12, "/", "", false, true)
	context.JSON(http.StatusOK, conf.SuccessResult(&LoginResult{
		UserId:    u.Id,
		AuthToken: "",
		Expires:   1,
	}))

}

//根据登录账号查找
func findUser(username string) *db.User {
	ok, err := regexp.Match(`^1\d{10}$`, []byte(username))
	if err != nil {
		panic(err)
	}
	if ok {
		return db.UserDao.GetByPhone(username)
	}
	if strings.Index(username, "@") != -1 {
		return db.UserDao.GetByEmail(username)
	}
	return nil
}

func CurrentUser(ctx *gin.Context) {
	val, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusOK, conf.ErrorResult(conf.LoginErrorStatus, "用户未登录"))
		return
	}
	u := val.(db.User)
	userVO := &UserVO{
		Userid:    u.Id,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Signature: u.Signature,
		Title:     u.Title,
		Address:   u.Address,
		Phone:     u.Phone,
		Access:    u.RoleNo,
	}
	ctx.JSON(http.StatusOK, conf.SuccessResult(userVO))
}

func OutLogin(ctx *gin.Context) {
	ctx.SetCookie(CookieName, "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, conf.SuccessResult(true))
}
