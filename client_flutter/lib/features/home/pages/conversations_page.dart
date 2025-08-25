import 'dart:async';
import 'package:flutter/material.dart';
import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/ui/we_colors.dart';
import 'package:client_flutter/core/util/base_list_page_state.dart';

class ConversationsPage extends StatefulWidget {
  const ConversationsPage({super.key});

  @override
  State<ConversationsPage> createState() => _ConversationsPageState();
}

class _ConversationsPageState
    extends BaseListPageState<Map<String, dynamic>, ConversationsPage> {
  final _chat = ChatService();
  StreamSubscription<Map<String, dynamic>>? _sub;

  @override
  Future<List<Map<String, dynamic>>> fetch({
    required int page,
    required int pageSize,
  }) async {
    final data = await _chat.getConversations(page: page, pageSize: pageSize);
    final raw = (data['items'] as List?)?.cast<Map>() ?? [];
    return raw.cast<Map<String, dynamic>>();
  }

  void _openChat(Map<String, dynamic> item) {
    if (item['type'] == 'group') {
      final g = item['group'] as Map?;
      if (g == null) return;
      Navigator.of(context)
          .pushNamed(
            '/chat',
            arguments: {
              'group_id': g['id'],
              'id': g['id'],
              'is_group': true,
              'nickname': g['name'] ?? '群聊',
              'avatar': g['avatar'] ?? '',
            },
          )
          .then((_) => reload());
      return;
    }
    final peerId = item['peer_id'];
    final peer = item['peer'] as Map?; // {id,nickname,avatar}
    Navigator.of(context)
        .pushNamed(
          '/chat',
          arguments: {
            'id': peerId,
            'nickname':
                (peer != null &&
                        (peer['nickname'] as String?)?.isNotEmpty == true)
                    ? peer['nickname']
                    : '用户$peerId',
            'email': '',
            'avatar': peer != null ? (peer['avatar'] ?? '') : '',
          },
        )
        .then((_) => reload());
  }

  @override
  @override
  Widget buildItem(BuildContext context, Map<String, dynamic> it, int index) {
    final last = (it['last_message'] as Map?)?.cast<String, dynamic>();
    final unread = (it['unread_count'] ?? 0) as int;
    String? avatar;
    String nickname;
    if (it['type'] == 'group') {
      final g = it['group'] as Map?;
      avatar = g != null ? g['avatar'] as String? : null;
      nickname =
          (g != null && (g['name'] as String?)?.isNotEmpty == true)
              ? g['name']
              : '群聊';
    } else {
      final peerId = it['peer_id'];
      final peer = it['peer'] as Map?; // {id,nickname,avatar}
      avatar = peer != null ? peer['avatar'] as String? : null;
      nickname =
          (peer != null && (peer['nickname'] as String?)?.isNotEmpty == true)
              ? peer['nickname']
              : '用户$peerId';
    }
    // 预览文案：图片/视频/链接使用占位符，不展示原始 URL
    String preview = '';
    if (last != null) {
      final type = (last['type'] ?? '').toString().toLowerCase();
      final content = (last['content'] ?? '').toString();
      if (type == 'image') {
        preview = '[图片]';
      } else if (type == 'video') {
        preview = '[视频]';
      } else if (type == 'link') {
        // 若后端未来提供显式 link type
        preview = '[链接]';
      } else if (content.startsWith('http://') ||
          content.startsWith('https://')) {
        // 文本消息里若是 URL，也做链接占位
        preview = '[链接]';
      } else {
        preview = content;
      }
    }
    return WeCell(
      leading: CircleAvatar(
        backgroundColor: Colors.grey.shade200,
        backgroundImage:
            (avatar != null && avatar.isNotEmpty) ? NetworkImage(avatar) : null,
        child:
            (avatar == null || avatar.isEmpty)
                ? const Icon(Icons.person_outline, color: Colors.white)
                : null,
      ),
      title: nickname,
      subtitle: preview,
      trailing:
          unread > 0
              ? Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: WeColors.badge,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  unread > 99 ? '99+' : unread.toString(),
                  style: const TextStyle(color: Colors.white, fontSize: 10),
                ),
              )
              : const SizedBox.shrink(),
      onTap: () => _openChat(it),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '小绿书',
        actions: [
          IconButton(onPressed: _showAddMenu, icon: const Icon(Icons.add)),
        ],
      ),
      body: super.build(context),
    );
  }

  void _showAddMenu() async {
    final RenderBox box = context.findRenderObject() as RenderBox;
    final offset = box.localToGlobal(Offset.zero);
    final selected = await showMenu<String>(
      context: context,
      position: RelativeRect.fromLTRB(
        offset.dx + box.size.width - 40,
        offset.dy + kToolbarHeight + 8,
        16,
        0,
      ),
      items: const [
        PopupMenuItem(value: 'group', child: Text('发起群聊')),
        PopupMenuItem(value: 'add_friend', child: Text('添加朋友')),
      ],
    );
    if (selected == 'group') {
      Navigator.of(context).pushNamed('/create_group');
    } else if (selected == 'add_friend') {
      Navigator.of(
        context,
      ).pushNamed('/contacts'); // TODO: implement friend add flow
    }
  }

  @override
  void initState() {
    super.initState();
    _sub = _chat.messageStream.listen((msg) {
      final sid = msg['sender_id'];
      final rid = msg['receiver_id'];
      if (sid == null || rid == null) return;
      bool matched = false;
      for (final it in items) {
        final pid = it['peer_id'];
        if (pid == sid || pid == rid) {
          it['last_message'] = msg;
          if (pid == sid) {
            // 对方发来的消息，增加未读
            final unread = (it['unread_count'] ?? 0) as int;
            it['unread_count'] = unread + 1;
          }
          matched = true;
        }
      }
      if (matched) {
        setState(() {
          items.sort((a, b) {
            final aid = (a['last_message']?['id'] ?? 0) as int;
            final bid = (b['last_message']?['id'] ?? 0) as int;
            return bid.compareTo(aid);
          });
        });
      } else {
        // 新会话，刷新获取最新列表
        reload();
      }
    });
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }
}
