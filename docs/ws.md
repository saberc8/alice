我来补一份“私聊即时通讯（WebSocket）使用说明”，涵盖鉴权、连接、发消息与拉历史的最小闭环。

## 快速开始

- 后端默认地址: http://localhost:8090
- API BasePath: /api/v1
- 鉴权: Bearer Token（移动端 App JWT）
- 前置条件:
  - 已通过 `/api/v1/app/register` 和 `/api/v1/app/login` 获取 `token`
  - 双方已互为好友（否则发送消息会报错 not friends）

## WebSocket 连接

- 连接地址: `ws://<host>:<port>/api/v1/app/chat/ws`
- 鉴权方式（两选一）
  - Header: `Authorization: Bearer <token>`
  - 或 Query: `?token=<token>`（浏览器直接用 WebSocket 场景更方便）

示例（浏览器/前端 JavaScript）:
- URL: `ws://localhost:8090/api/v1/app/chat/ws?token=YOUR_APP_TOKEN`

示例（Node.js，带 Header）:
- 头部使用 `Authorization: Bearer <token>` 连接

## 消息发送与接收

- 客户端发送 JSON（单条）:
  - 字段: 
    - type: string，可选，默认 "text"
    - to: number，必填，对端用户的 app_user_id
    - content: string，必填，消息内容
- 成功后：
  - 服务端会将消息持久化，并回写给发送方；
  - 若接收方在线，也会实时推送相同的消息对象。

发送示例（客户端 -> 服务端）:
{
  "type": "text",
  "to": 1024,
  "content": "hello"
}

接收示例（服务端 -> 客户端）:
{
  "id": 123,
  "sender_id": 1001,
  "receiver_id": 1024,
  "type": "text",
  "content": "hello",
  "is_read": false,
  "read_at": null,
  "created_at": "2025-08-14T08:00:00Z"
}

说明:
- 仅支持点对点消息（单聊）。
- 必须是好友关系才允许发送；否则返回错误 not friends。
- 当前消息类型支持 text（可扩展）。

## 历史消息查询（REST）

- 路径: `GET /api/v1/app/chat/history/{peer_id}`
- 鉴权: Bearer（App Token）
- 查询参数:
  - page: 默认 1
  - page_size: 默认 20，最大 100
- 响应:
  - data.items: Message[]（按时间倒序）
  - data.total: 总数
  - data.page/page_size

响应示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 123,
        "sender_id": 1001,
        "receiver_id": 1024,
        "type": "text",
        "content": "hello",
        "is_read": false,
        "read_at": null,
        "created_at": "2025-08-14T08:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}

用法提示:
- peer_id 为对端 app_user_id。
- 历史记录包含双向消息（我->对方、对方->我）。

## 鉴权与错误

- 未带 token 或 token 无效: 401 unauthorized
- 非好友发送: 400 not friends
- 参数不合法（如 content 为空、to 自己）: 400 invalid params

## 表结构与持久化

- 表名: `app_chat_messages`
- 字段: id, sender_id, receiver_id, type, content, is_read, read_at, created_at
- 服务端在发送时写库；历史接口从库分页查询。

## 开发联调要点

- 本地启动: go run main.go
- WebSocket 可用 Query 方式带 token，便于浏览器直连调试。
- 先通过好友请求/接受接口建立好友关系，再测试发消息。
- Swagger（若启用）: http://localhost:8090/swagger/index.html 可查看接口概览（BasePath: /api/v1）。

需要补充的增强项（可选）:
- 送达/已读回执（MarkRead REST 与 WS 推送）
- 离线消息未读数统计
- 富媒体消息类型（图片/文件）与上传策略
- 会话列表与最后一条消息/未读统计聚合

如需，我可以补一个 Postman/HTTPie/前端 demo 脚本，或加“已读上报”接口与推送。
