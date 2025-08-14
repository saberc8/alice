package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"alice/api/model"
	"alice/infra/config"
	"alice/pkg/logger"
)

var (
	ErrTokenInvalid = errors.New("invalid token")
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
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

		// 设置用户ID到上下文
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "invalid token payload"))
			c.Abort()
			return
		}

		c.Set("user_id", uint(userID))
		c.Next()
	}
}

// extractToken 从请求中提取token
func extractToken(c *gin.Context) string {
	// 从Header中获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 从查询参数中获取
	if t := c.Query("token"); t != "" {
		return t
	}
	// 兼容常见参数名
	if t := c.Query("access_token"); t != "" {
		return t
	}
	return ""
}

// validateToken 验证token
func validateToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, ErrTokenInvalid
	}

	uid, ok := userID.(uint)
	if !ok {
		return 0, ErrTokenInvalid
	}

	return uid, nil
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		logger.Infof("%s %s %d", start, path, status)
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
