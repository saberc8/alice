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
    final peerId = item['peer_id'];
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
  @override
  Widget buildItem(BuildContext context, Map<String, dynamic> it, int index) {
    final last = (it['last_message'] as Map?)?.cast<String, dynamic>();
    final unread = (it['unread_count'] ?? 0) as int;
    final peerId = it['peer_id'];
    final preview = last != null ? (last['content']?.toString() ?? '') : '';
    return WeCell(
      leading: CircleAvatar(
        backgroundColor: Colors.grey.shade300,
        child: const Icon(Icons.person_outline, color: Colors.white),
      ),
      title: '用户$peerId',
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
        title: '微信',
        actions: [
          IconButton(onPressed: reload, icon: const Icon(Icons.refresh)),
        ],
      ),
      body: super.build(context),
    );
  }
}
