import 'dart:async';

import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:flutter/material.dart';
import 'package:client_flutter/ui/we_colors.dart';
import 'package:client_flutter/ui/widgets/emoji_picker.dart';
import 'package:image_picker/image_picker.dart';
import 'package:url_launcher/url_launcher.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key, required this.peer});
  final Map<String, dynamic> peer; // {id, email, nickname, avatar, bio}

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _svc = ChatService();
  final _msgCtrl = TextEditingController();
  final FocusNode _inputFocus = FocusNode();
  bool _showEmoji = false; // 是否显示表情面板

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
      final lastIncoming = _messages.where((m) => m['sender_id'] == _peerId);
      if (lastIncoming.isNotEmpty) {
        final bid = (lastIncoming.last['id'] as num?)?.toInt() ?? 0;
        if (bid > 0) unawaited(_svc.markRead(peerId: _peerId, beforeId: bid));
      }

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
            if (sid == _peerId) {
              final bid = (event['id'] as num?)?.toInt() ?? 0;
              if (bid > 0)
                unawaited(_svc.markRead(peerId: _peerId, beforeId: bid));
            }
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
    // 发送后保持输入焦点
    if (!_inputFocus.hasFocus) _inputFocus.requestFocus();
  }

  Future<void> _sendImage() async {
    try {
      final picker = ImagePicker();
      final img = await picker.pickImage(
        source: ImageSource.gallery,
        imageQuality: 85,
      );
      if (img == null) return;
      final url = await _svc.uploadImage(img.path);
      if (url == null || _sink == null) return;
      _sink!.add({'type': 'image', 'to': _peerId, 'content': url});
    } catch (e) {
      if (mounted)
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('图片发送失败: $e')));
    }
  }

  void _sendLink() async {
    final ctrl = TextEditingController();
    final ok = await showDialog<bool>(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text('发送链接'),
            content: TextField(
              controller: ctrl,
              decoration: const InputDecoration(
                hintText: '输入以 http/https 开头的链接',
              ),
              autofocus: true,
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text('取消'),
              ),
              FilledButton(
                onPressed: () => Navigator.pop(ctx, true),
                child: const Text('发送'),
              ),
            ],
          ),
    );
    if (ok != true) return;
    final url = ctrl.text.trim();
    if (url.isEmpty || _sink == null) return;
    // 简单校验
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      if (mounted)
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('链接格式不正确')));
      return;
    }
    _sink!.add({'type': 'link', 'to': _peerId, 'content': url});
  }

  void _toggleEmojiPanel() {
    // 如果键盘在，就先收起键盘再展示表情
    if (_showEmoji) {
      setState(() => _showEmoji = false);
      // 回到输入焦点
      _inputFocus.requestFocus();
    } else {
      // 取消文本焦点 -> 收起系统键盘
      _inputFocus.unfocus();
      setState(() => _showEmoji = true);
    }
  }

  void _onEmojiSelected(String emoji) {
    final text = _msgCtrl.text;
    final sel = _msgCtrl.selection;
    final insertAt = sel.isValid ? sel.start : text.length;
    final newText = text.replaceRange(insertAt, insertAt, emoji);
    _msgCtrl.value = TextEditingValue(
      text: newText,
      selection: TextSelection.collapsed(offset: insertAt + emoji.length),
    );
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
        centerTitle: true,
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
                          final type = m['type']?.toString() ?? 'text';
                          Widget contentWidget;
                          if (type == 'image') {
                            final url = m['content']?.toString() ?? '';
                            contentWidget = GestureDetector(
                              onTap: () async {
                                if (await canLaunchUrl(Uri.parse(url))) {
                                  await launchUrl(
                                    Uri.parse(url),
                                    mode: LaunchMode.externalApplication,
                                  );
                                }
                              },
                              child: ClipRRect(
                                borderRadius: BorderRadius.circular(8),
                                child: Image.network(
                                  url,
                                  width: 180,
                                  height: 180,
                                  fit: BoxFit.cover,
                                  errorBuilder:
                                      (_, __, ___) => const SizedBox(
                                        width: 120,
                                        height: 120,
                                        child: Center(
                                          child: Icon(Icons.broken_image),
                                        ),
                                      ),
                                ),
                              ),
                            );
                          } else if (type == 'link') {
                            final url = m['content']?.toString() ?? '';
                            contentWidget = InkWell(
                              onTap: () async {
                                if (await canLaunchUrl(Uri.parse(url))) {
                                  await launchUrl(
                                    Uri.parse(url),
                                    mode: LaunchMode.externalApplication,
                                  );
                                }
                              },
                              child: Text(
                                url,
                                style: const TextStyle(
                                  color: Colors.blue,
                                  decoration: TextDecoration.underline,
                                ),
                              ),
                            );
                          } else {
                            contentWidget = Text(
                              m['content']?.toString() ?? '',
                              style: TextStyle(
                                color: isMe ? Colors.black : Colors.black87,
                              ),
                            );
                          }
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
                                        ? WeColors.bubbleMe
                                        : WeColors.bubbleOther,
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: contentWidget,
                            ),
                          );
                        },
                      ),
                    ),
          ),
          SafeArea(
            top: false,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Row(
                  children: [
                    IconButton(
                      onPressed: _toggleEmojiPanel,
                      icon: Icon(
                        _showEmoji
                            ? Icons.keyboard_alt_outlined
                            : Icons.tag_faces_outlined,
                        color: WeColors.textSecondary,
                      ),
                    ),
                    IconButton(
                      onPressed: _sendImage,
                      icon: const Icon(
                        Icons.image_outlined,
                        color: WeColors.textSecondary,
                      ),
                    ),
                    IconButton(
                      onPressed: _sendLink,
                      icon: const Icon(
                        Icons.link,
                        color: WeColors.textSecondary,
                      ),
                    ),
                    Expanded(
                      child: TextField(
                        focusNode: _inputFocus,
                        controller: _msgCtrl,
                        decoration: const InputDecoration(
                          hintText: '发消息',
                          border: InputBorder.none,
                          contentPadding: EdgeInsets.symmetric(horizontal: 8),
                        ),
                        minLines: 1,
                        maxLines: 4,
                        onTap: () {
                          if (_showEmoji) setState(() => _showEmoji = false);
                        },
                      ),
                    ),
                    IconButton(
                      onPressed: _send,
                      icon: const Icon(
                        Icons.send_rounded,
                        color: WeColors.brand,
                      ),
                    ),
                  ],
                ),
                AnimatedSwitcher(
                  duration: const Duration(milliseconds: 200),
                  child:
                      _showEmoji
                          ? SizedBox(
                            key: const ValueKey('emoji'),
                            height: 300,
                            child: EmojiPicker(
                              onSelected: (e) => _onEmojiSelected(e),
                            ),
                          )
                          : const SizedBox.shrink(key: ValueKey('empty')),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
