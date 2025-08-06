/*
 * Copyright 2025 alice Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package router

import (
	"github.com/gin-gonic/gin"

	"alice/api/handler"
	"alice/api/middleware"
)

type Router struct {
	userHandler *handler.UserHandler
}

func NewRouter(userHandler *handler.UserHandler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := v1.Group("/users")
		{
			// 公开接口
			userGroup.POST("/register", r.userHandler.Register)
			userGroup.POST("/login", r.userHandler.Login)

			// 需要认证的接口
			authenticated := userGroup.Group("")
			authenticated.Use(middleware.JWTAuth())
			{
				authenticated.GET("/profile", r.userHandler.GetProfile)
				authenticated.PUT("/profile", r.userHandler.UpdateProfile)
			}
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
