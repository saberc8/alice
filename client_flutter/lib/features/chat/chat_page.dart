import 'dart:async';

import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:flutter/material.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key, required this.peer});
  final Map<String, dynamic> peer; // {id, email, nickname, avatar, bio}

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _svc = ChatService();
  final _msgCtrl = TextEditingController();

  final List<Map<String, dynamic>> _messages = [];
  StreamSubscription<Map<String, dynamic>>? _sub;
  StreamSink<Map<String, dynamic>>? _sink;
  Future<void> Function()? _close;
  bool _loading = true;
  String? _error;
  int _page = 1;
  bool _hasMore = true;

  int get _peerId {
    final v = widget.peer['id'];
    if (v is int) return v;
    if (v is double) return v.toInt();
    if (v is String) return int.tryParse(v) ?? 0;
    return 0;
  }

  String get _title {
    final name = (widget.peer['nickname'] as String?) ?? '';
    final email = (widget.peer['email'] as String?) ?? '';
    return name.isNotEmpty ? name : email;
  }

  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<void> _init() async {
    try {
      // history first
      final data = await _svc.getHistory(
        peerId: _peerId,
        page: 1,
        pageSize: 20,
      );
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      _messages.addAll(items.cast<Map<String, dynamic>>().reversed);

      final (stream, sink, close) = _svc.connect();
      _sink = sink;
      _close = close;
      _sub = stream.listen(
        (event) {
          // Only append messages that belong to this conversation
          final sid = event['sender_id'];
          final rid = event['receiver_id'];
          if (sid == _peerId || rid == _peerId) {
            setState(() => _messages.add(event));
          }
        },
        onError: (e) {
          setState(() => _error = e.toString());
        },
      );
    } catch (e) {
      setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _loadMore() async {
    if (!_hasMore) return;
    final next = _page + 1;
    try {
      final data = await _svc.getHistory(
        peerId: _peerId,
        page: next,
        pageSize: 20,
      );
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      if (items.isEmpty) {
        setState(() => _hasMore = false);
      } else {
        setState(() {
          _page = next;
          _messages.insertAll(0, items.cast<Map<String, dynamic>>().reversed);
        });
      }
    } catch (_) {}
  }

  void _send() {
    final text = _msgCtrl.text.trim();
    if (text.isEmpty || _sink == null) return;
    _sink!.add({'type': 'text', 'to': _peerId, 'content': text});
    _msgCtrl.clear();
  }

  @override
  void dispose() {
    _sub?.cancel();
    _close?.call();
    _msgCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final avatar = widget.peer['avatar'] as String?;

    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            CircleAvatar(
              radius: 16,
              backgroundImage:
                  (avatar != null && avatar.isNotEmpty)
                      ? NetworkImage(avatar)
                      : null,
              child:
                  (avatar == null || avatar.isEmpty)
                      ? const Icon(Icons.person_outline, size: 16)
                      : null,
            ),
            const SizedBox(width: 12),
            Text(_title),
          ],
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child:
                _loading
                    ? const Center(child: CircularProgressIndicator())
                    : _error != null
                    ? Center(child: Text('加载失败: $_error'))
                    : NotificationListener<ScrollNotification>(
                      onNotification: (n) {
                        if (n.metrics.pixels <= 40 &&
                            n is ScrollUpdateNotification) {
                          _loadMore();
                        }
                        return false;
                      },
                      child: ListView.builder(
                        padding: const EdgeInsets.symmetric(
                          vertical: 12,
                          horizontal: 12,
                        ),
                        itemCount: _messages.length,
                        itemBuilder: (ctx, i) {
                          final m = _messages[i];
                          final isMe =
                              m['receiver_id'] == _peerId ? true : false;
                          // In absence of my id, infer by receiver_id equals peer
                          return Align(
                            alignment:
                                isMe
                                    ? Alignment.centerRight
                                    : Alignment.centerLeft,
                            child: Container(
                              margin: const EdgeInsets.symmetric(vertical: 4),
                              padding: const EdgeInsets.symmetric(
                                vertical: 8,
                                horizontal: 12,
                              ),
                              decoration: BoxDecoration(
                                color:
                                    isMe
                                        ? Theme.of(context).colorScheme.primary
                                        : Colors.grey.shade200,
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Text(
                                m['content']?.toString() ?? '',
                                style: TextStyle(
                                  color: isMe ? Colors.white : Colors.black87,
                                ),
                              ),
                            ),
                          );
                        },
                      ),
                    ),
          ),
          SafeArea(
            top: false,
            child: Row(
              children: [
                IconButton(
                  onPressed: () {},
                  icon: const Icon(Icons.add_circle_outline),
                ),
                Expanded(
                  child: TextField(
                    controller: _msgCtrl,
                    decoration: const InputDecoration(
                      hintText: '发消息',
                      border: InputBorder.none,
                      contentPadding: EdgeInsets.symmetric(horizontal: 8),
                    ),
                    minLines: 1,
                    maxLines: 4,
                  ),
                ),
                IconButton(
                  onPressed: _send,
                  icon: const Icon(Icons.send_rounded),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
