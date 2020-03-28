package service

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"../PO"
	"../utils"

	jwt "github.com/appleboy/gin-jwt"
	goJwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	identifyKey = "jwtInfo"
	jwtSecKey   = `Y'/JZU\CLxdpeyMZeR\DvmlD2R0_zy@`
)

var (
	autnInrequirePathes                       = []string{}
	adminCheckPathes                          = []string{}
	jwtMW               *jwt.GinJWTMiddleware = nil
	cookieDomain                              = ""
)

type JwtContent struct {
	OpenID  string `json:"openId"`
	UserID  string `json:"userId"`
	AppName string `json:"appName"`
}

func UpdateCookieDomain(domain string) {
	if nil != jwtMW {
		jwtMW.CookieDomain = domain
	}
	cookieDomain = domain
}

func JwtParseContent(c *gin.Context) JwtContent {
	dic, _ := jwtMW.GetClaimsFromJWT(c)
	utils.DebugLogger.Logf("jwt parse result %v", dic)
	ret := JwtContent{}
	if _, ok := dic["openId"].(string); ok {
		ret.OpenID = dic["openId"].(string)
	}
	if _, ok := dic["userId"].(string); ok {
		ret.UserID = dic["userId"].(string)
	}
	return ret
}

func generateJwtToken(content JwtContent) (string, int, error) {
	token := goJwt.New(goJwt.GetSigningMethod(jwtMW.SigningAlgorithm))
	claims := token.Claims.(goJwt.MapClaims)
	expire := jwtMW.TimeFunc().Add(jwtMW.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = jwtMW.TimeFunc().Unix()
	claims["appName"] = content.AppName
	claims["userId"] = content.UserID
	claims["openId"] = content.OpenID

	t := reflect.TypeOf(content)
	v := reflect.ValueOf(content)
	i := 0
	for i < t.NumField() {
		f := t.Field(i)
		name := f.Name
		jsonTagStr := f.Tag.Get("json")
		if len(jsonTagStr) > 0 {
			name = strings.Split(jsonTagStr, ",")[0]
		}
		claims[name] = v.Field(i).Interface()
		i++
	}

	tokenString, err := token.SignedString(jwtMW.Key)
	if nil != err {
		return "", 0, err
	}
	maxage := int(expire.Unix() - time.Now().Unix())
	return tokenString, maxage, nil
}

func JwtLogin(ctx *gin.Context, content JwtContent) error {
	tokenString, maxage, err := generateJwtToken(content)
	if nil != err {
		return err
	}
	ctx.Header(jwtMW.TokenHeadName, tokenString)
	ctx.SetCookie(
		jwtMW.CookieName,
		tokenString,
		maxage,
		"/",
		jwtMW.CookieDomain,
		jwtMW.SecureCookie,
		jwtMW.CookieHTTPOnly,
	)
	return nil
}

func JwtLogout(ctx *ServiceCtx) {
	ctx.Ctx.SetCookie(
		jwtMW.CookieName,
		"",
		-1,
		"/",
		jwtMW.CookieDomain,
		jwtMW.SecureCookie,
		jwtMW.CookieHTTPOnly,
	)
}

func authMiddleWare() func(*gin.Context) {
	jwtConf, err := jwt.New(&jwt.GinJWTMiddleware{
		CookieDomain:   cookieDomain,
		Key:            []byte(jwtSecKey),
		SendCookie:     true,
		SecureCookie:   false,   //non HTTPS dev environments
		CookieHTTPOnly: true,    // JS can't modify
		CookieName:     "token", // default jwt
		TokenHeadName:  "Token",
		TokenLookup:    "cookie:token,header:Token",
		Timeout:        time.Hour * 24 * 180,
		Unauthorized: func(c *gin.Context, code int, message string) {
			utils.DebugLogger.Logf("auth fail")
			PO.Error_NeedLogin.RendFail(c)
		},
	})
	if nil != err {
		fmt.Println(err)
		panic(err)
	}
	jwtMW = jwtConf

	return func(c *gin.Context) {
		utils.DebugLogger.Logf("start auth")
		if !pathNeedAuth(c.Request.URL.Path) {
			utils.DebugLogger.Logf("no need auth")
			c.Next()
			return
		}

		if IsAdmin(c) {
			utils.DebugLogger.Logf("is admin")
			c.Next()
			return
		}

		if pathNeedAdmin(c.Request.URL.Path) {
		} else {
			utils.DebugLogger.Logf("normal check")
			claims, err := jwtMW.GetClaimsFromJWT(c)
			if nil == err && claims["exp"] != nil {
				if _, ok := claims["exp"].(float64); ok {
					if int64(claims["exp"].(float64)) > jwtMW.TimeFunc().Unix() {
						c.Set("JWT_PAYLOAD", claims)
						c.Next()
						return
					}
				}
			}
		}

		c.Abort()
		jwtMW.Unauthorized(c, 0, "")
	}
}

func pathNeedAuth(path string) bool {
	for _, v := range autnInrequirePathes {
		if v == path || strings.HasPrefix(path, v) {
			return false
		}
	}
	return true
}

func pathNeedAdmin(path string) bool {
	for _, v := range adminCheckPathes {
		if v == path {
			utils.DebugLogger.Logf("Need admin check %v", path)
			return true
		}
	}
	return false
}

func RegisterAuthNoNeedPath(path string) {
	autnInrequirePathes = append(autnInrequirePathes, path)
}

func JwtHeaderToken(ctx *gin.Context) string {
	return ctx.Writer.Header().Get(jwtMW.TokenHeadName)
}

func RegisterAdminCheckPath(path string) {
	adminCheckPathes = append(adminCheckPathes, path)
}

func IsAdmin(c *gin.Context) bool {
	utils.InfoLogger.Logf("admin check ...")
	adminUserIDs := []string{"cc3b199babd34718aa7578a80f6b0a06"}
	userId := JwtParseContent(c).UserID
	for i := 0; i < len(adminUserIDs); i++ {
		if userId == adminUserIDs[i] {
			utils.InfoLogger.Logf("admin check succeed")
			return true
		}
	}
	utils.InfoLogger.Logf("admin check fail")
	return false
}
