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

package entity

import (
	"time"
)

// UserRole 用户角色关联表
type UserRole struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   string    `json:"user_id" gorm:"not null;type:varchar(36);index"`
	RoleID   string    `json:"role_id" gorm:"not null;type:varchar(36);index"`
	CreateAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoleID       string    `json:"role_id" gorm:"not null;type:varchar(36);index"`
	PermissionID string    `json:"permission_id" gorm:"not null;type:varchar(36);index"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    string    `json:"role_id" gorm:"not null;type:varchar(36);index"`
	MenuID    string    `json:"menu_id" gorm:"not null;type:varchar(36);index"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}
