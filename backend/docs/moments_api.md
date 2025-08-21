# 朋友圈动态接口文档

版本: v1  基础路径: `/api/v1`

> 鉴权：除特别说明外均需要携带 App Token
>
> Header: `Authorization: Bearer <token>`

## 目录
1. [发布动态](#发布动态)
2. [全部动态列表](#全部动态列表)
3. [查看指定用户动态列表](#查看指定用户动态列表)
4. [删除动态](#删除动态)
5. [数据模型](#数据模型)
6. [错误码说明](#错误码说明)
7. [设计要点](#设计要点)

---
## 发布动态
POST `/api/v1/app/moments`

发布一条公开的朋友圈动态（当前所有用户可见）。

### 请求头
Authorization: Bearer <token>
Content-Type: application/json

### 请求体
```jsonc
{
  "content": "今天天气不错～",          // 必填，非空
  "images": ["/app-avatars/xxx.jpg", "/moments/img-2.png"] // 选填，最多 9 张，相对路径
}
```

规则：
- images 超过 9 条截断
- 空字符串自动过滤
- 不校验文件是否真实存在（前端应先上传得到路径）

### 成功响应 200
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 12,
    "user_id": 3,
    "nickname": "Alice",
    "avatar": "http://<minio-base>/app-avatars/avatar-3-1734923000.jpg",
    "content": "今天天气不错～",
    "images": ["/app-avatars/xxx.jpg", "/moments/img-2.png"],
    "created_at": 1734923010
  }
}
```

### 可能错误
| HTTP | code | message             | 场景 |
|------|------|---------------------|------|
| 400  | 400  | invalid request     | JSON 解析失败 / content 为空 |
| 401  | 401  | unauthorized        | Token 缺失或无效 |
| 500  | 500  | internal server error | 服务内部错误 |

---
## 全部动态列表
GET `/api/v1/app/moments`

按时间倒序返回所有用户的动态。

### 查询参数
| 名称 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int  | 否 | 默认 1 |
| page_size | int | 否 | 默认 20，最大 100 |

### 成功响应 200
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 12,
        "user_id": 3,
        "nickname": "Alice",
        "avatar": "http://<minio-base>/app-avatars/avatar-3-1734923000.jpg",
        "content": "今天天气不错～",
        "images": ["/app-avatars/xxx.jpg"],
        "created_at": 1734923010
      }
    ],
    "total": 35,
    "page": 1,
    "page_size": 20
  }
}
```

---
## 查看指定用户动态列表
GET `/api/v1/app/users/{user_id}/moments`

查看某个用户发布的动态。

路径参数：
- `user_id` 目标用户 ID

查询参数与响应结构同“全部动态列表”。

### 可能错误
| HTTP | code | message | 场景 |
|------|------|---------|------|
| 400  | 400  | invalid request | user_id 非法 |
| 401  | 401  | unauthorized   | 未认证 |
| 500  | 500  | internal server error | 内部错误 |

---
## 删除动态
DELETE `/api/v1/app/moments/{moment_id}`

删除当前登录用户自己发布的动态。

路径参数：
- `moment_id` 动态 ID

### 成功响应 200
```json
{
  "code": 200,
  "message": "deleted"
}
```

### 可能错误
| HTTP | code | message         | 场景 |
|------|------|-----------------|------|
| 400  | 400  | invalid request | 参数解析失败 |
| 400  | 400  | <错误信息>      | 业务校验失败（如未找到或已删除） |
| 401  | 401  | unauthorized    | 未认证 |

> 当前后端简单按 `id AND user_id` 条件删除，若行不存在返回影响 0 行（映射 400）。后续可细化为 404。

---
## 数据模型

### CreateMomentRequest
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 文本内容 |
| images | []string | 否 | 相对路径（/bucket/object），最多 9 个 |

### MomentItem
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 动态 ID |
| user_id | uint | 发布者 ID |
| nickname | string | 发布者昵称 |
| avatar | string | 发布者头像完整 URL |
| content | string | 内容 |
| images | []string | 图片相对路径数组 |
| created_at | int64 | 创建时间（Unix 秒） |

### MomentListResponse
| 字段 | 类型 | 说明 |
|------|------|------|
| items | []MomentItem | 动态数组 |
| total | int64 | 符合条件总数 |
| page | int | 当前页码 |
| page_size | int | 每页大小 |

---
## 错误码说明
| code | message | 解释 |
|------|---------|------|
| 200 | success | 成功 |
| 400 | invalid request / 具体业务错误 | 入参或业务校验错误 |
| 401 | unauthorized | 未认证或 token 失效 |
| 404 | user not found 等 | 资源未找到（当前动态删除未用） |
| 500 | internal server error | 服务端异常 |

---
## 设计要点
- 目前所有动态公开显示；后续可扩展可见范围（public / friends / private）。
- 图片路径统一使用“相对路径”保存，响应中头像会补全为完整可访问 URL；动态图片暂保持原样，前端可自行拼 baseURL（亦可后端补全，可后续调整）。
- 未实现软删除；如需要审计可增加 `deleted_at` 字段并改为软删除。
- Swagger 注释已存在（发布、列表、用户列表、删除），运行 `swag init` 可同步生成文档。
- 可扩展功能：点赞、评论、转发、屏蔽、仅好友可见、话题标签、全文搜索。

---
最后更新: 自动生成
