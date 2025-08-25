import 'package:flutter/foundation.dart';
import '../data/moment_api.dart';
import '../data/moment_models.dart';

class MomentStore extends ChangeNotifier {
  final MomentApi api;
  final int? userId; // 若提供则按用户过滤
  MomentStore({MomentApi? api, this.userId}) : api = api ?? MomentApi();

  final List<MomentItem> _moments = [];
  bool _isLoading = false;
  bool _hasMore = true;
  int _page = 1;
  final int _pageSize = 20;

  List<MomentItem> get moments => List.unmodifiable(_moments);
  bool get isLoading => _isLoading;
  bool get hasMore => _hasMore;

  Future<void> refresh() async {
    _page = 1;
    _hasMore = true;
    _moments.clear();
    await loadMore();
  }

  Future<void> loadMore() async {
    if (_isLoading || !_hasMore) return;
    _isLoading = true;
    notifyListeners();
    try {
      final res =
          userId == null
              ? await api.listAll(page: _page, pageSize: _pageSize)
              : await api.listByUser(userId!, page: _page, pageSize: _pageSize);
      if (_page == 1) _moments.clear();
      _moments.addAll(res.items);
      _hasMore = _moments.length < res.total;
      if (_hasMore) _page += 1;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<MomentItem?> publish(String content, List<String> images) async {
    final m = await api.publish(content: content, images: images);
    // 新发布的放到顶部
    _moments.insert(0, m);
    notifyListeners();
    return m;
  }

  Future<void> delete(int id) async {
    await api.delete(id);
    _moments.removeWhere((e) => e.id == id);
    notifyListeners();
  }

  Future<void> like(int id) async {
    final idx = _moments.indexWhere((e) => e.id == id);
    if (idx == -1) return;
    final m = _moments[idx];
    if (m.liked) return;
    _moments[idx] = MomentItem(
      id: m.id,
      userId: m.userId,
      nickname: m.nickname,
      avatar: m.avatar,
      content: m.content,
      images: m.images,
      createdAt: m.createdAt,
      likeCount: m.likeCount + 1,
      liked: true,
    );
    notifyListeners();
    try {
      await api.like(id);
    } catch (_) {
      // revert
      _moments[idx] = m;
      notifyListeners();
    }
  }

  Future<void> unlike(int id) async {
    final idx = _moments.indexWhere((e) => e.id == id);
    if (idx == -1) return;
    final m = _moments[idx];
    if (!m.liked) return;
    _moments[idx] = MomentItem(
      id: m.id,
      userId: m.userId,
      nickname: m.nickname,
      avatar: m.avatar,
      content: m.content,
      images: m.images,
      createdAt: m.createdAt,
      likeCount: (m.likeCount - 1).clamp(0, 1 << 31),
      liked: false,
    );
    notifyListeners();
    try {
      await api.unlike(id);
    } catch (_) {
      _moments[idx] = m;
      notifyListeners();
    }
  }
}
