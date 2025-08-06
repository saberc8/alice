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

// PermissionStatus 权限状态
type PermissionStatus string

const (
	PermissionStatusActive   PermissionStatus = "active"
	PermissionStatusInactive PermissionStatus = "inactive"
)

// Permission 权限实体
type Permission struct {
	ID          string           `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string           `json:"name" gorm:"not null;size:100"`
	Code        string           `json:"code" gorm:"uniqueIndex;not null;size:100"`
	Resource    string           `json:"resource" gorm:"not null;size:100"`
	Action      string           `json:"action" gorm:"not null;size:50"`
	Description *string          `json:"description" gorm:"size:500"`
	Status      PermissionStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// IsActive 检查权限是否激活
func (p *Permission) IsActive() bool {
	return p.Status == PermissionStatusActive
}
