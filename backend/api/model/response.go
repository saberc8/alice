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

package model

// APIResponse 标准API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterResponse 注册响应数据
type RegisterResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 响应状态码常量
const (
	CodeSuccess         = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeInternalError   = 500
	CodeValidationError = 422
)

// 响应消息常量
const (
	MsgSuccess            = "success"
	MsgLoginSuccess       = "login successful"
	MsgRegisterSuccess    = "register successful"
	MsgLogoutSuccess      = "logout successful"
	MsgInvalidRequest     = "invalid request"
	MsgInvalidCredentials = "invalid username or password"
	MsgUnauthorized       = "unauthorized"
	MsgUserNotFound       = "user not found"
	MsgUserAlreadyExists  = "user already exists"
	MsgInternalError      = "internal server error"
)

// 辅助函数创建标准响应
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	}
}

func SuccessResponseWithMessage(message string, data interface{}) APIResponse {
	return APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
	}
}
