import 'package:flutter/material.dart';
import 'package:client_flutter/core/app/profile_service.dart';
import 'package:client_flutter/core/app/friends_service.dart';
import 'package:client_flutter/features/contacts/friend_profile_page.dart';
import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:client_flutter/ui/we_tabbar.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/ui/we_colors.dart';

class HomeTabs extends StatefulWidget {
  const HomeTabs({super.key, required this.onLogout});

  final VoidCallback onLogout;

  @override
  State<HomeTabs> createState() => _HomeTabsState();
}

class _HomeTabsState extends State<HomeTabs> {
  int _index = 0;

  final _pages = const [
    _ConversationsPage(),
    _ContactsPage(),
    _DiscoverPage(),
    _ProfilePage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(index: _index, children: _pages),
      bottomNavigationBar: WeTabBar(
        items: const [
          WeTabItem(
            icon: Icons.chat_bubble_outline,
            iconActive: Icons.chat_bubble,
            label: '微信',
          ),
          WeTabItem(
            icon: Icons.contacts_outlined,
            iconActive: Icons.contacts,
            label: '通讯录',
          ),
          WeTabItem(
            icon: Icons.explore_outlined,
            iconActive: Icons.explore,
            label: '发现',
          ),
          WeTabItem(
            icon: Icons.person_outline,
            iconActive: Icons.person,
            label: '我',
          ),
        ],
        currentIndex: _index,
        onTap: (i) => setState(() => _index = i),
      ),
    );
  }
}

class _ConversationsPage extends StatefulWidget {
  const _ConversationsPage();
  @override
  State<_ConversationsPage> createState() => _ConversationsPageState();
}

class _ConversationsPageState extends State<_ConversationsPage> {
  final _chat = ChatService();
  List<Map<String, dynamic>> _items = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final data = await _chat.getConversations();
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      if (!mounted) return;
      setState(() => _items = items.cast<Map<String, dynamic>>());
    } catch (e) {
      if (!mounted) return;
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _openChat(Map<String, dynamic> item) {
    final peerId = item['peer_id'];
    // 简单推断 peer 基本字段（实际可请求用户资料）
    Navigator.of(context).pushNamed(
      '/chat',
      arguments: {
        'id': peerId,
        'nickname': '用户$peerId',
        'email': '',
        'avatar': '',
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '微信',
        actions: [
          IconButton(onPressed: _load, icon: const Icon(Icons.refresh)),
        ],
      ),
      body:
          _loading
              ? const Center(child: CircularProgressIndicator())
              : _error != null
              ? Center(child: Text('加载失败: $_error'))
              : RefreshIndicator(
                onRefresh: _load,
                child:
                    _items.isEmpty
                        ? ListView(
                          children: const [
                            SizedBox(height: 160),
                            Center(child: Text('暂无会话')),
                          ],
                        )
                        : ListView.separated(
                          itemCount: _items.length,
                          separatorBuilder:
                              (_, __) => const Divider(
                                height: 1,
                                color: WeColors.divider,
                              ),
                          itemBuilder: (ctx, i) {
                            final it = _items[i];
                            final last =
                                (it['last_message'] as Map?)
                                    ?.cast<String, dynamic>();
                            final unread = (it['unread_count'] ?? 0) as int;
                            final peerId = it['peer_id'];
                            final preview =
                                last != null
                                    ? (last['content']?.toString() ?? '')
                                    : '';
                            return WeCell(
                              leading: CircleAvatar(
                                backgroundColor: Colors.grey.shade300,
                                child: const Icon(
                                  Icons.person_outline,
                                  color: Colors.white,
                                ),
                              ),
                              title: '用户$peerId',
                              subtitle: preview,
                              trailing:
                                  unread > 0
                                      ? Container(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 6,
                                          vertical: 2,
                                        ),
                                        decoration: BoxDecoration(
                                          color: WeColors.badge,
                                          borderRadius: BorderRadius.circular(
                                            12,
                                          ),
                                        ),
                                        child: Text(
                                          unread > 99
                                              ? '99+'
                                              : unread.toString(),
                                          style: const TextStyle(
                                            color: Colors.white,
                                            fontSize: 10,
                                          ),
                                        ),
                                      )
                                      : const SizedBox.shrink(),
                              onTap: () => _openChat(it),
                            );
                          },
                        ),
              ),
    );
  }
}

class _ContactsPage extends StatefulWidget {
  const _ContactsPage();
  @override
  State<_ContactsPage> createState() => _ContactsPageState();
}

class _ContactsPageState extends State<_ContactsPage> {
  final _svc = FriendsService();
  List<Map<String, dynamic>> _friends = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final data = await _svc.getFriends();
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      if (!mounted) return;
      setState(() => _friends = items.cast<Map<String, dynamic>>());
    } catch (e) {
      if (!mounted) return;
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _addFriendDialog() async {
    final ctrl = TextEditingController();
    final email = await showDialog<String>(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text('添加好友'),
            content: TextField(
              controller: ctrl,
              autofocus: true,
              decoration: const InputDecoration(hintText: '输入好友邮箱'),
              keyboardType: TextInputType.emailAddress,
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(ctx).pop(),
                child: const Text('取消'),
              ),
              FilledButton(
                onPressed: () => Navigator.of(ctx).pop(ctrl.text),
                child: const Text('发送请求'),
              ),
            ],
          ),
    );
    if (email == null || email.trim().isEmpty) return;
    try {
      await _svc.sendFriendRequest(email);
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('好友请求已发送')));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('发送失败: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '通讯录',
        actions: [
          IconButton(
            onPressed:
                () => Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) => const _FriendRequestsPage(),
                  ),
                ),
            icon: const Icon(Icons.inbox_outlined),
            tooltip: '待处理请求',
          ),
          IconButton(
            onPressed: _addFriendDialog,
            icon: const Icon(Icons.person_add_alt_1),
            tooltip: '添加好友',
          ),
        ],
      ),
      body:
          _loading
              ? const Center(child: CircularProgressIndicator())
              : _error != null
              ? Center(child: Text('加载失败: $_error'))
              : RefreshIndicator(
                onRefresh: _load,
                child:
                    _friends.isEmpty
                        ? ListView(
                          children: const [
                            SizedBox(height: 160),
                            Center(child: Text('暂无好友')),
                          ],
                        )
                        : ListView.separated(
                          itemCount: _friends.length,
                          separatorBuilder:
                              (_, __) => const Divider(
                                height: 1,
                                color: WeColors.divider,
                              ),
                          itemBuilder: (ctx, i) {
                            final u = _friends[i];
                            return WeCell(
                              leading: CircleAvatar(
                                backgroundImage:
                                    (u['avatar'] != null &&
                                            (u['avatar'] as String).isNotEmpty)
                                        ? NetworkImage(u['avatar'])
                                        : null,
                                child:
                                    (u['avatar'] == null ||
                                            (u['avatar'] as String).isEmpty)
                                        ? const Icon(Icons.person)
                                        : null,
                              ),
                              title:
                                  u['nickname']?.toString().isNotEmpty == true
                                      ? u['nickname']
                                      : (u['email'] ?? '-'),
                              subtitle: u['email'] ?? '',
                              onTap: () {
                                Navigator.of(context).push(
                                  MaterialPageRoute(
                                    builder: (_) => FriendProfilePage(user: u),
                                  ),
                                );
                              },
                            );
                          },
                        ),
              ),
    );
  }
}

class _FriendRequestsPage extends StatefulWidget {
  const _FriendRequestsPage();
  @override
  State<_FriendRequestsPage> createState() => _FriendRequestsPageState();
}

class _FriendRequestsPageState extends State<_FriendRequestsPage> {
  final _svc = FriendsService();
  List<int> _requestIds = [];
  List<int> _requesterIds = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final data = await _svc.getPendingRequests();
      if (!mounted) return;
      setState(() {
        _requestIds = (data['request_ids'] as List?)?.cast<int>() ?? [];
        _requesterIds = (data['requester_ids'] as List?)?.cast<int>() ?? [];
      });
    } catch (e) {
      if (!mounted) return;
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _accept(int index) async {
    final id = _requestIds[index];
    try {
      await _svc.acceptRequest(id);
      await _load();
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('已接受')));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('操作失败: $e')));
    }
  }

  Future<void> _decline(int index) async {
    final id = _requestIds[index];
    try {
      await _svc.declineRequest(id);
      await _load();
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('已拒绝')));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('操作失败: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const WeAppBar(title: '待处理好友请求'),
      body:
          _loading
              ? const Center(child: CircularProgressIndicator())
              : _error != null
              ? Center(child: Text('加载失败: $_error'))
              : RefreshIndicator(
                onRefresh: _load,
                child:
                    (_requestIds.isEmpty)
                        ? ListView(
                          children: const [
                            SizedBox(height: 160),
                            Center(child: Text('暂无待处理请求')),
                          ],
                        )
                        : ListView.separated(
                          itemCount: _requestIds.length,
                          separatorBuilder:
                              (_, __) => const Divider(
                                height: 1,
                                color: WeColors.divider,
                              ),
                          itemBuilder: (ctx, i) {
                            final reqId = _requestIds[i];
                            final requesterId =
                                i < _requesterIds.length
                                    ? _requesterIds[i]
                                    : null;
                            return WeCell(
                              leading: const CircleAvatar(
                                child: Icon(Icons.person_outline),
                              ),
                              title: '请求者 ID: ${requesterId ?? '-'}',
                              subtitle: '请求 ID: $reqId',
                              trailing: Wrap(
                                spacing: 8,
                                children: [
                                  TextButton(
                                    onPressed: () => _decline(i),
                                    child: const Text('拒绝'),
                                  ),
                                  FilledButton(
                                    onPressed: () => _accept(i),
                                    child: const Text('接受'),
                                  ),
                                ],
                              ),
                            );
                          },
                        ),
              ),
    );
  }
}

class _DiscoverPage extends StatelessWidget {
  const _DiscoverPage();
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: WeAppBar(title: '发现'),
      body: Center(child: Text('发现内容占位')),
    );
  }
}

class _ProfilePage extends StatefulWidget {
  const _ProfilePage();
  @override
  State<_ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<_ProfilePage> {
  final _svc = ProfileService();
  Map<String, dynamic>? _profile;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final data = await _svc.getProfile();
      if (!mounted) return;
      setState(() => _profile = data);
    } catch (e) {
      if (!mounted) return;
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  void _logout() {
    final state = context.findAncestorStateOfType<_HomeTabsState>();
    state?.widget.onLogout();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '我',
        actions: [
          IconButton(onPressed: _load, icon: const Icon(Icons.refresh)),
          IconButton(onPressed: _logout, icon: const Icon(Icons.logout)),
        ],
      ),
      body:
          _loading
              ? const Center(child: CircularProgressIndicator())
              : _error != null
              ? Center(child: Text('加载失败: $_error'))
              : _profile == null
              ? const Center(child: Text('暂无资料'))
              : Column(
                children: [
                  const SizedBox(height: 12),
                  Material(
                    color: Colors.white,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                      child: Row(
                        children: [
                          CircleAvatar(
                            radius: 28,
                            backgroundImage:
                                (_profile!['avatar'] != null &&
                                        (_profile!['avatar'] as String)
                                            .isNotEmpty)
                                    ? NetworkImage(_profile!['avatar'])
                                    : null,
                            child:
                                (_profile!['avatar'] == null ||
                                        (_profile!['avatar'] as String).isEmpty)
                                    ? const Icon(Icons.person, size: 28)
                                    : null,
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  _profile!['nickname'] ?? '-',
                                  style: const TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  _profile!['email'] ?? '-',
                                  style: const TextStyle(
                                    color: WeColors.textSecondary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const Icon(Icons.qr_code, color: Colors.grey),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Expanded(
                    child: ListView(
                      children: const [
                        WeCell(
                          title: '支付',
                          leading: Icon(
                            Icons.payment_outlined,
                            color: Colors.black87,
                          ),
                        ),
                        Divider(height: 1, color: WeColors.divider),
                        WeCell(
                          title: '收藏',
                          leading: Icon(
                            Icons.star_outline,
                            color: Colors.black87,
                          ),
                        ),
                        Divider(height: 1, color: WeColors.divider),
                        WeCell(
                          title: '设置',
                          leading: Icon(
                            Icons.settings_outlined,
                            color: Colors.black87,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
    );
  }
}
