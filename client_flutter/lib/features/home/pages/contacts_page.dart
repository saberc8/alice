import 'package:flutter/material.dart';
import 'package:client_flutter/core/app/friends_service.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/features/contacts/friend_profile_page.dart';
import 'friend_requests_page.dart';
import 'package:client_flutter/core/util/base_list_page_state.dart';

class ContactsPage extends StatefulWidget {
  const ContactsPage({super.key});
  @override
  State<ContactsPage> createState() => _ContactsPageState();
}

class _ContactsPageState
    extends BaseListPageState<Map<String, dynamic>, ContactsPage> {
  final _svc = FriendsService();

  @override
  Future<List<Map<String, dynamic>>> fetch({
    required int page,
    required int pageSize,
  }) async {
    final data = await _svc.getFriends(page: page, pageSize: pageSize);
    final items = (data['items'] as List?)?.cast<Map>() ?? [];
    return items.cast<Map<String, dynamic>>();
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
  @override
  Widget buildItem(BuildContext context, Map<String, dynamic> u, int index) {
    return WeCell(
      leading: CircleAvatar(
        backgroundImage:
            (u['avatar'] != null && (u['avatar'] as String).isNotEmpty)
                ? NetworkImage(u['avatar'])
                : null,
        child:
            (u['avatar'] == null || (u['avatar'] as String).isEmpty)
                ? const Icon(Icons.person)
                : null,
      ),
      title:
          u['nickname']?.toString().isNotEmpty == true
              ? u['nickname']
              : (u['email'] ?? '-'),
      subtitle: u['email'] ?? '',
      onTap: () {
        Navigator.of(
          context,
        ).push(MaterialPageRoute(builder: (_) => FriendProfilePage(user: u)));
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '通讯录',
        actions: [
          IconButton(
            onPressed: reload,
            icon: const Icon(Icons.refresh),
            tooltip: '刷新',
          ),
          IconButton(
            onPressed:
                () => Navigator.of(context).push(
                  MaterialPageRoute(builder: (_) => const FriendRequestsPage()),
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
      body: super.build(context),
    );
  }
}
