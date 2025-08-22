import 'package:flutter/material.dart';

/// 抽象封装通用列表页的加载/错误/刷新逻辑，降低重复样板代码。
/// 子类需实现 [fetch] 拉取数据 与 [buildItem] 渲染单行。
abstract class BaseListPageState<T, W extends StatefulWidget> extends State<W> {
  final List<T> items = [];
  bool loading = true;
  String? error;

  /// 分页参数（如后续需要可扩展），当前简单保留
  int page = 1;
  int pageSize = 20;
  bool hasMore = true;
  bool _loadingMore = false;

  /// 拉取首页数据（必须实现）
  Future<List<T>> fetch({required int page, required int pageSize});

  /// 可选：自定义空视图
  Widget buildEmpty(BuildContext context) => Column(
    mainAxisSize: MainAxisSize.min,
    children: [
      const Icon(Icons.inbox_outlined, size: 56, color: Colors.grey),
      const SizedBox(height: 12),
      const Text('暂无数据', style: TextStyle(color: Colors.grey)),
    ],
  );

  /// 构建单个列表项
  Widget buildItem(BuildContext context, T item, int index);

  /// 可选：分割线
  Widget buildSeparator(BuildContext context, int index) =>
      const Divider(height: 1);

  @override
  void initState() {
    super.initState();
    reload();
  }

  Future<void> reload() async {
    setState(() {
      loading = true;
      error = null;
      page = 1;
      hasMore = true;
      items.clear();
    });
    try {
      final data = await fetch(page: page, pageSize: pageSize);
      if (!mounted) return;
      setState(() {
        items.addAll(data);
        hasMore = data.length == pageSize; // 简单判断是否还有更多
      });
    } catch (e) {
      if (!mounted) return;
      setState(() => error = e.toString());
      _toast('加载失败: $e');
    } finally {
      if (mounted) setState(() => loading = false);
    }
  }

  Future<void> loadMore() async {
    if (!hasMore || _loadingMore) return;
    _loadingMore = true;
    try {
      final nextPage = page + 1;
      final data = await fetch(page: nextPage, pageSize: pageSize);
      if (!mounted) return;
      setState(() {
        page = nextPage;
        items.addAll(data);
        hasMore = data.length == pageSize;
      });
    } catch (e) {
      debugPrint('Load more failed: $e');
      _toast('加载更多失败');
    } finally {
      _loadingMore = false;
    }
  }

  /// 可选：底部加载更多指示器
  Widget? buildFooter(BuildContext context) {
    if (!hasMore) return const SizedBox.shrink();
    return const Padding(
      padding: EdgeInsets.symmetric(vertical: 16),
      child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return _buildSkeleton();
    }
    if (error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('加载失败: $error'),
            const SizedBox(height: 12),
            FilledButton(onPressed: reload, child: const Text('重试')),
          ],
        ),
      );
    }
    if (items.isEmpty) {
      return RefreshIndicator(
        onRefresh: reload,
        child: ListView(
          children: [
            const SizedBox(height: 160),
            Center(child: buildEmpty(context)),
          ],
        ),
      );
    }
    return RefreshIndicator(
      onRefresh: reload,
      child: NotificationListener<ScrollNotification>(
        onNotification: (n) {
          if (n.metrics.pixels >= n.metrics.maxScrollExtent - 100) {
            loadMore();
          }
          return false;
        },
        child: ListView.separated(
          itemCount: items.length + 1,
          separatorBuilder:
              (c, i) =>
                  i < items.length - 1
                      ? buildSeparator(c, i)
                      : const SizedBox.shrink(),
          itemBuilder: (c, i) {
            if (i == items.length) {
              return buildFooter(c) ?? const SizedBox.shrink();
            }
            return buildItem(c, items[i], i);
          },
        ),
      ),
    );
  }

  void _toast(String msg) {
    final messenger = ScaffoldMessenger.maybeOf(context);
    if (messenger != null) {
      messenger.showSnackBar(SnackBar(content: Text(msg)));
    }
  }

  Widget _buildSkeleton() {
    return ListView.separated(
      physics: const NeverScrollableScrollPhysics(),
      padding: const EdgeInsets.symmetric(vertical: 12),
      itemCount: 8,
      separatorBuilder: (_, __) => const Divider(height: 1),
      itemBuilder: (c, i) => const _SkeletonTile(),
    );
  }
}

class _SkeletonTile extends StatelessWidget {
  const _SkeletonTile();
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: Colors.grey.shade300,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _line(width: 140),
                const SizedBox(height: 8),
                _line(width: 200),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _line({required double width}) => Container(
    width: width,
    height: 12,
    decoration: BoxDecoration(
      color: Colors.grey.shade300,
      borderRadius: BorderRadius.circular(4),
    ),
  );
}
