import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../state/moment_store.dart';
import '../data/moment_models.dart';
import 'publish_moment_sheet.dart';
import 'package:client_flutter/core/network/dio_client.dart';

class MomentListPage extends StatelessWidget {
  const MomentListPage({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => MomentStore()..refresh(),
      child: const _MomentListBody(),
    );
  }
}

class _MomentListBody extends StatelessWidget {
  const _MomentListBody();
  @override
  Widget build(BuildContext context) {
    final store = context.watch<MomentStore>();
    return Scaffold(
      appBar: AppBar(
        title: const Text('朋友圈'),
        actions: [
          IconButton(
            tooltip: '发布',
            icon: const Icon(Icons.create),
            onPressed: () async {
              final store = context.read<MomentStore>();
              await showModalBottomSheet(
                context: context,
                isScrollControlled: true,
                builder:
                    (_) => ChangeNotifierProvider.value(
                      value: store,
                      child: const PublishMomentSheet(),
                    ),
              );
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: store.refresh,
        child: NotificationListener<ScrollNotification>(
          onNotification: (n) {
            if (n.metrics.pixels >= n.metrics.maxScrollExtent - 200 &&
                store.hasMore &&
                !store.isLoading) {
              store.loadMore();
            }
            return false;
          },
          child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            itemCount: store.moments.length + 1,
            separatorBuilder: (_, __) => const Divider(height: 0),
            itemBuilder: (context, index) {
              if (index == store.moments.length) {
                if (store.isLoading) {
                  return const Padding(
                    padding: EdgeInsets.all(16),
                    child: Center(child: CircularProgressIndicator()),
                  );
                }
                if (!store.hasMore) {
                  return const Padding(
                    padding: EdgeInsets.all(16),
                    child: Center(
                      child: Text(
                        '没有更多了',
                        style: TextStyle(color: Colors.grey),
                      ),
                    ),
                  );
                }
                return const SizedBox.shrink();
              }
              final m = store.moments[index];
              return _MomentItemTile(item: m);
            },
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          final store = context.read<MomentStore>();
          await showModalBottomSheet(
            context: context,
            isScrollControlled: true,
            builder:
                (_) => ChangeNotifierProvider.value(
                  value: store,
                  child: const PublishMomentSheet(),
                ),
          );
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}

class _MomentItemTile extends StatelessWidget {
  final MomentItem item;
  const _MomentItemTile({required this.item});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              CircleAvatar(
                backgroundImage: NetworkImage(item.avatar),
                radius: 20,
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  item.nickname,
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
              ),
              Text(
                _formatTime(item.createdAt),
                style: const TextStyle(color: Colors.grey, fontSize: 12),
              ),
              const SizedBox(width: 4),
              PopupMenuButton<String>(
                onSelected: (v) async {
                  if (v == 'delete') {
                    final confirmed = await showDialog<bool>(
                      context: context,
                      builder:
                          (_) => AlertDialog(
                            title: const Text('删除动态'),
                            content: const Text('确定删除这条动态吗？'),
                            actions: [
                              TextButton(
                                onPressed: () => Navigator.pop(context, false),
                                child: const Text('取消'),
                              ),
                              ElevatedButton(
                                onPressed: () => Navigator.pop(context, true),
                                child: const Text('删除'),
                              ),
                            ],
                          ),
                    );
                    if (confirmed == true) {
                      await context.read<MomentStore>().delete(item.id);
                    }
                  }
                },
                itemBuilder:
                    (_) => const [
                      PopupMenuItem(value: 'delete', child: Text('删除')),
                    ],
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(item.content),
          if (item.images.isNotEmpty) ...[
            const SizedBox(height: 8),
            _MomentImagesGrid(images: item.images),
          ],
        ],
      ),
    );
  }
}

class _MomentImagesGrid extends StatelessWidget {
  final List<String> images;
  const _MomentImagesGrid({required this.images});
  @override
  Widget build(BuildContext context) {
    final count = images.length;
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: count == 1 ? 1 : (count <= 4 ? 2 : 3),
        crossAxisSpacing: 4,
        mainAxisSpacing: 4,
      ),
      itemCount: count,
      itemBuilder:
          (_, i) => ClipRRect(
            borderRadius: BorderRadius.circular(6),
            child: Image.network(_resolveImage(images[i]), fit: BoxFit.cover),
          ),
    );
  }
}

String _formatTime(DateTime dt) {
  final now = DateTime.now();
  final diff = now.difference(dt);
  if (diff.inMinutes < 1) return '刚刚';
  if (diff.inMinutes < 60) return '${diff.inMinutes}分钟前';
  if (diff.inHours < 24) return '${diff.inHours}小时前';
  return '${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}';
}

String _resolveImage(String raw) {
  if (raw.startsWith('http://') || raw.startsWith('https://')) return raw;
  // 后端返回相对路径时自行拼 base（此处简单用同域 API 基础，生产可单独配置）
  return '${DioClient().dio.options.baseUrl}$raw';
}
