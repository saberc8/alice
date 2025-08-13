package middleware

import (
	"net/http"
	"strconv"

	"alice/api/model"
	"alice/domain/rbac/service"
	"alice/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限中间件
func PermissionMiddleware(permissionService service.PermissionService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中获取用户ID (这里假设JWT中间件已经处理并设置了用户ID)
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(http.StatusUnauthorized, "用户未认证"))
			c.Abort()
			return
		}

		// 将用户ID转换为字符串
		var userID string
		switch v := userIDValue.(type) {
		case string:
			userID = v
		case uint:
			userID = strconv.FormatUint(uint64(v), 10)
		case int:
			userID = strconv.Itoa(v)
		default:
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(http.StatusUnauthorized, "无效的用户ID"))
			c.Abort()
			return
		}

		// 检查用户权限
		hasPermission, err := permissionService.CheckUserPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			logger.Errorf("检查用户权限失败: %v", err)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, "权限检查失败"))
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, model.ErrorResponse(http.StatusForbidden, "权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 辅助函数，创建权限中间件
func RequirePermission(permissionService service.PermissionService, resource, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, resource, action)
}
