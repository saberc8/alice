# Moments Feature (朋友圈)

包含：
- 数据模型 `MomentItem`
- 接口封装 `MomentApi`
- 状态管理 `MomentStore` (ChangeNotifier)
- UI 页面 `MomentListPage` + 发布弹出层 `PublishMomentSheet`

接入步骤：
1. 确保后端已运行，并提供 /api/v1/app/moments 接口。
2. 在登录后入口添加导航到 `MomentListPage()`。
3. 已使用 provider 包；若全局已有 MultiProvider，可将 `MomentStore` 提升至更高层。

图片上传：当前示例未真正上传，只是使用占位相对路径 `/mock/<filename>`；正式环境：
1. 先调用对象存储上传 API 获得相对路径（如 /moments/xxx.jpg）。
2. 调用 publish 传这些路径。

后续可扩展：
- 点赞 / 评论
- 仅好友可见
- 懒加载图片及缓存
- 详情页 / 下拉刷新动画优化
