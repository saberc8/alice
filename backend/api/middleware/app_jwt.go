package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"alice/api/model"
	"alice/infra/config"
)

// AppJWTAuth 专用于移动端用户的 JWT 中间件
func AppJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "missing authorization token"))
			c.Abort()
			return
		}
		cfg := config.Load()
		claims, err := validateToken(token, cfg.JWT.SecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "invalid token"))
			c.Abort()
			return
		}
		// 仅接受 app_user_id
		if id, ok := claims["app_user_id"].(float64); ok {
			c.Set("app_user_id", uint(id))
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "invalid token payload"))
		c.Abort()
	}
}

// 可与通用 validateToken 复用（同包）
func ValidateAppToken(tokenString string) (jwt.MapClaims, error) {
	cfg := config.Load()
	return validateToken(tokenString, cfg.JWT.SecretKey)
}
