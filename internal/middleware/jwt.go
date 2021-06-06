package middleware

import (
	"github.com/Elingering/go-web/global"
	"github.com/Elingering/go-web/internal/model"
	"github.com/Elingering/go-web/pkg/app"
	"github.com/Elingering/go-web/pkg/errcode"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}
		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			claims, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			} else {
				jwtAuth := model.Auth{
					AppKey:    claims.AppKey,
					AppSecret: claims.AppSecret,
				}
				auth, err := jwtAuth.Get(global.DBEngine)
				if err != nil {
					response := app.NewResponse(c)
					response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
				}
				c.Set("auth", auth)
				//auth, _ := c.Get("auth")//auth: auth.(model.Auth)
			}
		}

		if ecode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
