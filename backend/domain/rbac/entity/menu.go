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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// MenuType 菜单类型
type MenuType int

const (
	MenuTypeGroup     MenuType = 0 // 分组
	MenuTypeCatalogue MenuType = 1 // 目录
	MenuTypeMenu      MenuType = 2 // 菜单
	MenuTypeButton    MenuType = 3 // 按钮
)

// MenuStatus 菜单状态
type MenuStatus string

const (
	MenuStatusActive   MenuStatus = "active"
	MenuStatusInactive MenuStatus = "inactive"
)

// MenuMeta 菜单元数据
type MenuMeta struct {
	Icon         *string `json:"icon,omitempty"`
	Caption      *string `json:"caption,omitempty"`
	Info         *string `json:"info,omitempty"`
	Disabled     *bool   `json:"disabled,omitempty"`
	Auth         *bool   `json:"auth,omitempty"`
	Hidden       *bool   `json:"hidden,omitempty"`
	ExternalLink *string `json:"external_link,omitempty"`
	Component    *string `json:"component,omitempty"`
}

// Value 实现 driver.Valuer 接口，用于将 MenuMeta 转换为数据库值
func (m MenuMeta) Value() (driver.Value, error) {
	// 检查是否为空的MenuMeta
	if m == (MenuMeta{}) ||
		(m.Icon == nil && m.Caption == nil && m.Info == nil &&
			m.Disabled == nil && m.Auth == nil && m.Hidden == nil &&
			m.ExternalLink == nil && m.Component == nil) {
		return "{}", nil // 返回空的JSON对象而不是null
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口，用于从数据库值转换为 MenuMeta
func (m *MenuMeta) Scan(value interface{}) error {
	if value == nil {
		*m = MenuMeta{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into MenuMeta", value)
	}

	// 如果是空字节或空字符串，设置为空的MenuMeta
	if len(bytes) == 0 || string(bytes) == "" || string(bytes) == "null" {
		*m = MenuMeta{}
		return nil
	}

	// 尝试解析JSON
	if err := json.Unmarshal(bytes, m); err != nil {
		// 如果解析失败，记录错误并设置为空的MenuMeta
		fmt.Printf("Failed to unmarshal MenuMeta: %v, data: %s\n", err, string(bytes))
		*m = MenuMeta{}
		return nil // 不返回错误，避免中断整个查询
	}

	return nil
}

// Menu 菜单实体
type Menu struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ParentID    *string    `json:"parent_id" gorm:"type:varchar(36);index"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Code        string     `json:"code" gorm:"uniqueIndex;not null;size:100"`
	Path        *string    `json:"path" gorm:"size:200"`
	Type        MenuType   `json:"type" gorm:"not null;default:2"`
	Order       int        `json:"order" gorm:"default:0"`
	Status      MenuStatus `json:"status" gorm:"not null;default:'active'"`
	Meta        MenuMeta   `json:"meta" gorm:"type:json"`
	Description *string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联关系
	Children []*Menu `json:"children,omitempty" gorm:"-"`
	Parent   *Menu   `json:"parent,omitempty" gorm:"-"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// IsActive 检查菜单是否激活
func (m *Menu) IsActive() bool {
	return m.Status == MenuStatusActive
}

// IsGroup 检查是否为分组
func (m *Menu) IsGroup() bool {
	return m.Type == MenuTypeGroup
}

// IsCatalogue 检查是否为目录
func (m *Menu) IsCatalogue() bool {
	return m.Type == MenuTypeCatalogue
}

// IsMenu 检查是否为菜单
func (m *Menu) IsMenu() bool {
	return m.Type == MenuTypeMenu
}

// IsButton 检查是否为按钮
func (m *Menu) IsButton() bool {
	return m.Type == MenuTypeButton
}
