package middleware

import (
	"errors"

	"github.com/Acyclonepl/Blog-basedon-gin/pkg/app"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			_, err := app.ParseToken(token)
			if err != nil {
				var ve *jwt.ValidationError
				if errors.As(err, &ve) {
					switch ve.Errors {
					case jwt.ValidationErrorExpired:
						ecode = errcode.UnauthorizedTokenTimeout
					default:
						ecode = errcode.UnauthorizedTokenError
					}
				} else {
					// 非验证类错误（如解析失败）
					ecode = errcode.UnauthorizedTokenError
				}
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
