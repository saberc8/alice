import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../state/moment_store.dart';
import '../data/moment_models.dart';
import '../data/moment_api.dart';
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
        centerTitle: true,
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
          const SizedBox(height: 8),
          Row(
            children: [
              InkWell(
                onTap: () {
                  final store = context.read<MomentStore>();
                  if (item.liked) {
                    store.unlike(item.id);
                  } else {
                    store.like(item.id);
                  }
                },
                child: Row(
                  children: [
                    Icon(
                      item.liked ? Icons.favorite : Icons.favorite_border,
                      color: item.liked ? Colors.red : Colors.grey,
                      size: 20,
                    ),
                    const SizedBox(width: 4),
                    Text('${item.likeCount}'),
                  ],
                ),
              ),
              const SizedBox(width: 24),
              InkWell(
                onTap: () {
                  _showComments(context, item);
                },
                child: const Row(
                  children: [
                    Icon(Icons.comment, size: 20, color: Colors.grey),
                    SizedBox(width: 4),
                    Text('评论'),
                  ],
                ),
              ),
            ],
          ),
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

void _showComments(BuildContext context, MomentItem item) {
  showModalBottomSheet(
    context: context,
    isScrollControlled: true,
    builder: (_) => _CommentsSheet(item: item),
  );
}

class _CommentsSheet extends StatefulWidget {
  final MomentItem item;
  const _CommentsSheet({required this.item});
  @override
  State<_CommentsSheet> createState() => _CommentsSheetState();
}

class _CommentsSheetState extends State<_CommentsSheet> {
  final _controller = TextEditingController();
  final List<MomentCommentItem> _comments = [];
  bool _loading = false;
  bool _posting = false;
  int _page = 1;
  bool _hasMore = true;

  final _api = MomentApi();

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    if (_loading || !_hasMore) return;
    setState(() {
      _loading = true;
    });
    final res = await _api.listComments(
      widget.item.id,
      page: _page,
      pageSize: 50,
    );
    setState(() {
      _comments.addAll(res.items);
      _hasMore = _comments.length < res.total;
      if (_hasMore) _page += 1;
      _loading = false;
    });
  }

  Future<void> _submit() async {
    final text = _controller.text.trim();
    if (text.isEmpty || _posting) return;
    setState(() {
      _posting = true;
    });
    try {
      final cmt = await _api.addComment(widget.item.id, text);
      setState(() {
        _comments.insert(0, cmt);
        _controller.clear();
      });
    } finally {
      setState(() {
        _posting = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottom = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.only(bottom: bottom),
      child: SafeArea(
        top: false,
        child: SizedBox(
          height: MediaQuery.of(context).size.height * 0.6,
          child: Column(
            children: [
              const SizedBox(height: 8),
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 8),
              Expanded(
                child: NotificationListener<ScrollNotification>(
                  onNotification: (n) {
                    if (n.metrics.pixels >= n.metrics.maxScrollExtent - 50) {
                      _load();
                    }
                    return false;
                  },
                  child: ListView.builder(
                    itemCount: _comments.length + (_loading ? 1 : 0),
                    itemBuilder: (context, i) {
                      if (i >= _comments.length) {
                        return const Padding(
                          padding: EdgeInsets.all(12),
                          child: Center(child: CircularProgressIndicator()),
                        );
                      }
                      final c = _comments[i];
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundImage: NetworkImage(c.avatar),
                        ),
                        title: Text(c.nickname),
                        subtitle: Text(c.content),
                        trailing: Text(
                          _formatTime(c.createdAt),
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.grey,
                          ),
                        ),
                      );
                    },
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _controller,
                        decoration: const InputDecoration(
                          hintText: '说点什么...',
                          isDense: true,
                          border: OutlineInputBorder(),
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    _posting
                        ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                        : IconButton(
                          onPressed: _submit,
                          icon: const Icon(Icons.send),
                        ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
