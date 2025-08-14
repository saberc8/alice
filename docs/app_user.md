我来整理“注册/登录、加好友、好友列表（含资料）”的后端接口文档，重点标注邮箱小写规范化与鉴权使用。

## 目标与范围
- 用户注册/登录（邮箱为唯一账号，服务端自动转小写）
- 发送好友请求（通过对方邮箱，需对方确认）
- 查看待处理好友请求、接受/拒绝
- 获取好友列表（返回好友详细资料）

## 通用信息
- Base URL: /api/v1
- 鉴权方式: Bearer Token（仅注册/登录无需鉴权，其他 App 接口需要）
  - 请求头: Authorization: Bearer <token>
- 邮箱规范化: 服务端在注册、登录与添加好友时会将邮箱转为小写并去除首尾空格
- Swagger 文档: 若已在配置中开启，可访问 /swagger/index.html（BasePath 为 /api/v1）

## 接口一览

### 1) 注册
- 方法与路径: POST /api/v1/app/register
- 请求体:
  - email string 必填，合法邮箱
  - password string 必填，至少 6 位
  - nickname string 可选，最长 30
- 成功响应: 200
  - data: 
    - user: { id, email, nickname, avatar, bio }
    - token: string（App JWT）
- 失败:
  - 400 参数校验失败或邮箱已存在

说明:
- 邮箱作为唯一账号，服务端会将 email 统一转为小写存储
- 注册成功自动下发 token

### 2) 登录
- 方法与路径: POST /api/v1/app/login
- 请求体:
  - email string 必填
  - password string 必填
- 成功响应: 200
  - data: { token }
- 失败:
  - 400 参数校验失败
  - 401 账号或密码错误/账号不可用

说明:
- 登录时邮箱同样会被转为小写再匹配

### 3) 获取当前用户资料
- 方法与路径: GET /api/v1/app/profile
- 鉴权: Bearer
- 成功响应: 200
  - data: { id, email, nickname, avatar, bio }
- 失败:
  - 401 未认证
  - 404 用户不存在

### 4) 更新当前用户资料
- 方法与路径: PUT /api/v1/app/profile
- 鉴权: Bearer
- 请求体:
  - nickname string 可选，最长 30
  - avatar string 可选，URL
  - bio string 可选，最长 160
- 成功响应: 200
  - data: { id, email, nickname, avatar, bio }
- 失败:
  - 400 参数校验失败
  - 401 未认证

### 5) 发送好友请求（通过邮箱）
- 方法与路径: POST /api/v1/app/friends/request
- 鉴权: Bearer
- 请求体:
  - friend_email string 必填，对方邮箱
- 成功响应: 200
  - data: null
  - message: "request sent"
- 失败:
  - 400 参数错误或业务错误（例如用户不存在等）
  - 401 未认证

说明:
- 服务端会将 friend_email 转为小写后查询
- 当前实现对“已是好友”未做强约束；对“重复的待处理请求”是幂等的（同一对 requester→addressee 的 pending 不会重复创建）

### 6) 查看待处理好友请求
- 方法与路径: GET /api/v1/app/friends/requests
- 鉴权: Bearer
- 查询参数:
  - page int 可选，默认 1
  - page_size int 可选，默认 20，最大 100
- 成功响应: 200
  - data: {
      request_ids: number[],
      requester_ids: number[],
      total: number,
      page: number,
      page_size: number
    }
- 失败:
  - 401 未认证

说明:
- 列出“我作为收件人(addressee)”的挂起请求列表及其请求者 ID

### 7) 接受好友请求
- 方法与路径: POST /api/v1/app/friends/requests/{request_id}/accept
- 鉴权: Bearer
- 路径参数:
  - request_id int 必填
- 成功响应: 200
  - message: "accepted"
- 失败:
  - 400 请求无效或无权限（仅收件人可接受）
  - 401 未认证

说明:
- 接受后会建立双向好友关系（存两条 A→B、B→A）

### 8) 拒绝好友请求
- 方法与路径: POST /api/v1/app/friends/requests/{request_id}/decline
- 鉴权: Bearer
- 路径参数:
  - request_id int 必填
- 成功响应: 200
  - message: "declined"
- 失败:
  - 400 请求无效（所有权校验基础版）
  - 401 未认证

### 9) 好友列表（返回详细资料）
- 方法与路径: GET /api/v1/app/friends
- 鉴权: Bearer
- 查询参数:
  - page int 可选，默认 1
  - page_size int 可选，默认 20，最大 100
- 成功响应: 200
  - data: {
      items: [
        { id, email, nickname, avatar, bio },
        ...
      ],
      total: number,
      page: number,
      page_size: number
    }
- 失败:
  - 401 未认证

说明:
- 服务端先查好友 ID，再批量拉取用户资料后返回
