import 'dart:async';

import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:flutter/material.dart';
import 'package:client_flutter/ui/we_colors.dart';
import 'package:client_flutter/ui/widgets/emoji_picker.dart';
import 'package:image_picker/image_picker.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:video_player/video_player.dart';
import 'package:client_flutter/features/moments/ui/user_moment_list_page.dart';
import 'package:client_flutter/features/chat/group_manage_page.dart';

class ChatPage extends StatefulWidget {
  const ChatPage({super.key, required this.peer});
  final Map<String, dynamic>
  peer; // {id, email, nickname, avatar, bio, group_id?, is_group?}

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _svc = ChatService();
  final _msgCtrl = TextEditingController();
  final FocusNode _inputFocus = FocusNode();
  final ScrollController _scrollCtrl = ScrollController();
  bool _showEmoji = false; // 是否显示表情面板

  final List<Map<String, dynamic>> _messages = [];
  StreamSubscription<Map<String, dynamic>>? _sub;
  StreamSink<Map<String, dynamic>>? _sink;
  Future<void> Function()? _close;
  bool _loading = true;
  String? _error;
  int _page = 1;
  bool _hasMore = true;
  int? _selfId;

  bool get _isGroup =>
      widget.peer['is_group'] == true || widget.peer['group_id'] != null;

  int get _peerId {
    // For private chat, original peer id
    final v = widget.peer['id'];
    if (v is int) return v;
    if (v is double) return v.toInt();
    if (v is String) return int.tryParse(v) ?? 0;
    return 0;
  }

  int get _groupId {
    final v = widget.peer['group_id'];
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
      Map<String, dynamic> data;
      if (_isGroup) {
        data = await _svc.getGroupHistory(
          groupId: _groupId,
          page: 1,
          pageSize: 20,
        );
      } else {
        data = await _svc.getHistory(peerId: _peerId, page: 1, pageSize: 20);
      }
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      _messages.addAll(items.cast<Map<String, dynamic>>().reversed);
      _selfId = await _svc.selfId();
      if (!_isGroup) {
        final lastIncoming = _messages.where((m) => m['sender_id'] == _peerId);
        if (lastIncoming.isNotEmpty) {
          final bid = (lastIncoming.last['id'] as num?)?.toInt() ?? 0;
          if (bid > 0) unawaited(_svc.markRead(peerId: _peerId, beforeId: bid));
        }
      } else {
        // 群聊：进入时上报已读到当前最后一条
        final last = _messages.isNotEmpty ? _messages.last : null;
        final lastId = (last?['id'] as num?)?.toInt() ?? 0;
        if (lastId > 0)
          unawaited(_svc.markGroupRead(groupId: _groupId, beforeMsgId: lastId));
      }

      final (stream, sink, close) = _svc.connect();
      _sink = sink;
      _close = close;
      _sub = stream.listen(
        (event) {
          // Only append messages that belong to this conversation
          if (_isGroup) {
            final gid = event['group_id'];
            if (gid != null && gid == _groupId) {
              setState(() => _messages.add(event));
              // 收到群消息如果是他人发的并且在底部，推进已读
              final sid = (event['sender_id'] as num?)?.toInt();
              if (sid != null && sid != _selfId) {
                final mid = (event['message_id'] ?? event['id']) as num?;
                if (mid != null) {
                  unawaited(
                    _svc.markGroupRead(
                      groupId: _groupId,
                      beforeMsgId: mid.toInt(),
                    ),
                  );
                }
              }
            }
            return;
          }
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
      // 首帧渲染后滚动到底部
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => _scrollToBottom(jump: true),
      );
    }
  }

  Future<void> _loadMore() async {
    if (!_hasMore) return;
    final next = _page + 1;
    try {
      final prevMaxExtent =
          _scrollCtrl.hasClients ? _scrollCtrl.position.maxScrollExtent : null;
      final data =
          _isGroup
              ? await _svc.getGroupHistory(
                groupId: _groupId,
                page: next,
                pageSize: 20,
              )
              : await _svc.getHistory(
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
        // 维持视口位置不跳动
        if (prevMaxExtent != null) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!_scrollCtrl.hasClients) return;
            final newMax = _scrollCtrl.position.maxScrollExtent;
            final delta = newMax - prevMaxExtent;
            final newOffset = _scrollCtrl.offset + delta;
            _scrollCtrl.jumpTo(newOffset);
          });
        }
      }
    } catch (_) {}
  }

  void _sendTextByIME() {
    final text = _msgCtrl.text.trim();
    if (text.isEmpty || _sink == null) return;
    if (_isGroup) {
      _sink!.add({'type': 'text', 'group_id': _groupId, 'content': text});
    } else {
      _sink!.add({'type': 'text', 'to': _peerId, 'content': text});
    }
    _msgCtrl.clear();
    // 发送后稍后滚动到底部
    WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToBottom());
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
      if (_isGroup) {
        _sink!.add({'type': 'image', 'group_id': _groupId, 'content': url});
      } else {
        _sink!.add({'type': 'image', 'to': _peerId, 'content': url});
      }
    } catch (e) {
      if (mounted)
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('图片发送失败: $e')));
    }
  }

  Future<void> _sendVideo() async {
    try {
      final picker = ImagePicker();
      final vid = await picker.pickVideo(
        source: ImageSource.gallery,
        maxDuration: const Duration(minutes: 5),
      );
      if (vid == null) return;
      final url = await _svc.uploadVideo(vid.path);
      if (url == null || _sink == null) return;
      if (_isGroup) {
        _sink!.add({'type': 'video', 'group_id': _groupId, 'content': url});
      } else {
        _sink!.add({'type': 'video', 'to': _peerId, 'content': url});
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('视频发送失败: $e')));
      }
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
    if (_isGroup) {
      _sink!.add({'type': 'link', 'group_id': _groupId, 'content': url});
    } else {
      _sink!.add({'type': 'link', 'to': _peerId, 'content': url});
    }
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

  void _showMoreActions() {
    showModalBottomSheet(
      context: context,
      builder:
          (ctx) => SafeArea(
            child: SizedBox(
              height: 160,
              child: GridView.count(
                crossAxisCount: 4,
                children: [
                  _ActionIcon(
                    icon: Icons.image_outlined,
                    label: '图片',
                    onTap: () {
                      Navigator.pop(ctx);
                      _sendImage();
                    },
                  ),
                  _ActionIcon(
                    icon: Icons.videocam_outlined,
                    label: '视频',
                    onTap: () {
                      Navigator.pop(ctx);
                      _sendVideo();
                    },
                  ),
                  _ActionIcon(
                    icon: Icons.link,
                    label: '链接',
                    onTap: () {
                      Navigator.pop(ctx);
                      _sendLink();
                    },
                  ),
                ],
              ),
            ),
          ),
    );
  }

  @override
  void dispose() {
    _sub?.cancel();
    _close?.call();
    _msgCtrl.dispose();
    _scrollCtrl.dispose();
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
        actions: [
          if (_isGroup)
            IconButton(
              icon: const Icon(Icons.group_outlined),
              tooltip: '群管理',
              onPressed: () async {
                // 跳转群管理页（后续实现详细页面）
                final group = widget.peer;
                await Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (_) => GroupManagePage(group: group),
                  ),
                );
                // 返回后可选择刷新群成员/资料
                setState(() {});
              },
            ),
        ],
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
                        controller: _scrollCtrl,
                        padding: const EdgeInsets.symmetric(
                          vertical: 12,
                          horizontal: 12,
                        ),
                        itemCount: _messages.length,
                        itemBuilder: (ctx, i) {
                          final m = _messages[i];
                          // 判定是否我发出：优先 sender_id == selfId
                          final senderId = (m['sender_id'] as num?)?.toInt();
                          final isMe = senderId != null && senderId == _selfId;
                          // In absence of my id, infer by receiver_id equals peer
                          final type = m['type']?.toString() ?? 'text';
                          final sender =
                              (m['sender'] as Map?)?.cast<String, dynamic>();
                          final avatarUrl = sender?['avatar'] as String? ?? '';
                          final nickname = sender?['nickname'] as String? ?? '';
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
                          } else if (type == 'video') {
                            final url = m['content']?.toString() ?? '';
                            contentWidget = _VideoBubble(url: url);
                          } else {
                            contentWidget = Text(
                              m['content']?.toString() ?? '',
                              style: TextStyle(
                                color: isMe ? Colors.black : Colors.black87,
                              ),
                            );
                          }
                          return Padding(
                            padding: const EdgeInsets.symmetric(vertical: 4),
                            child: Row(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              mainAxisAlignment:
                                  isMe
                                      ? MainAxisAlignment.end
                                      : MainAxisAlignment.start,
                              children: [
                                if (!isMe)
                                  _AvatarButton(url: avatarUrl, user: sender),
                                if (!isMe) const SizedBox(width: 6),
                                Flexible(
                                  child: Column(
                                    crossAxisAlignment:
                                        isMe
                                            ? CrossAxisAlignment.end
                                            : CrossAxisAlignment.start,
                                    children: [
                                      if (_isGroup &&
                                          !isMe &&
                                          nickname.isNotEmpty)
                                        Padding(
                                          padding: const EdgeInsets.only(
                                            left: 4,
                                            right: 4,
                                            bottom: 2,
                                          ),
                                          child: Text(
                                            nickname,
                                            style: const TextStyle(
                                              fontSize: 11,
                                              color: Colors.black45,
                                            ),
                                          ),
                                        ),
                                      Container(
                                        padding: const EdgeInsets.symmetric(
                                          vertical: 8,
                                          horizontal: 12,
                                        ),
                                        decoration: BoxDecoration(
                                          color:
                                              isMe
                                                  ? WeColors.bubbleMe
                                                  : WeColors.bubbleOther,
                                          borderRadius: BorderRadius.circular(
                                            12,
                                          ),
                                        ),
                                        child: contentWidget,
                                      ),
                                    ],
                                  ),
                                ),
                                if (isMe) const SizedBox(width: 6),
                                if (isMe)
                                  _AvatarButton(url: avatarUrl, user: sender),
                              ],
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
                        textInputAction: TextInputAction.send,
                        onSubmitted: (_) => _sendTextByIME(),
                        onTap: () {
                          if (_showEmoji) setState(() => _showEmoji = false);
                        },
                      ),
                    ),
                    IconButton(
                      onPressed: _showMoreActions,
                      icon: const Icon(
                        Icons.add_circle_outline,
                        color: WeColors.textSecondary,
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

  bool _isNearBottom() {
    if (!_scrollCtrl.hasClients) return true;
    final pos = _scrollCtrl.position;
    return (pos.maxScrollExtent - pos.pixels) < 120;
  }

  void _scrollToBottom({bool jump = false}) {
    if (!_scrollCtrl.hasClients) return;
    if (jump) {
      _scrollCtrl.jumpTo(_scrollCtrl.position.maxScrollExtent);
    } else if (_isNearBottom()) {
      _scrollCtrl.animateTo(
        _scrollCtrl.position.maxScrollExtent,
        duration: const Duration(milliseconds: 250),
        curve: Curves.easeOut,
      );
    }
  }
}

class _ActionIcon extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  const _ActionIcon({
    required this.icon,
    required this.label,
    required this.onTap,
  });
  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          CircleAvatar(radius: 26, child: Icon(icon, size: 26)),
          const SizedBox(height: 6),
          Text(label, style: const TextStyle(fontSize: 12)),
        ],
      ),
    );
  }
}

class _AvatarButton extends StatelessWidget {
  final String url;
  final Map<String, dynamic>? user;
  const _AvatarButton({required this.url, required this.user});
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        if (user == null) return;
        Navigator.push(
          context,
          MaterialPageRoute(builder: (_) => UserMomentListPage(user: user!)),
        );
      },
      child: CircleAvatar(
        radius: 18,
        backgroundImage: (url.isNotEmpty) ? NetworkImage(url) : null,
        child: url.isEmpty ? const Icon(Icons.person_outline, size: 18) : null,
      ),
    );
  }
}

// 简单的视频气泡：首次点击进入播放模式，使用 VideoPlayer 控制器
class _VideoBubble extends StatefulWidget {
  const _VideoBubble({required this.url});
  final String url;
  @override
  State<_VideoBubble> createState() => _VideoBubbleState();
}

class _VideoBubbleState extends State<_VideoBubble> {
  VideoPlayerController? _controller;
  bool _initing = false;
  bool _playing = false;

  Future<void> _initAndPlay() async {
    if (_controller != null) {
      setState(() => _playing = !_playing);
      if (_playing) {
        _controller!.play();
      } else {
        _controller!.pause();
      }
      return;
    }
    setState(() => _initing = true);
    final c = VideoPlayerController.networkUrl(Uri.parse(widget.url));
    try {
      await c.initialize();
      c.setLooping(true);
      setState(() {
        _controller = c;
        _playing = true;
      });
      c.play();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('视频加载失败: $e')));
      }
      await c.dispose();
    } finally {
      if (mounted) setState(() => _initing = false);
    }
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final ctrl = _controller;
    final size = const Size(180, 240);
    Widget child;
    if (ctrl == null) {
      child = Stack(
        children: [
          Container(
            width: size.width,
            height: size.height,
            decoration: BoxDecoration(
              color: Colors.black12,
              borderRadius: BorderRadius.circular(8),
            ),
            child: Center(
              child:
                  _initing
                      ? const CircularProgressIndicator(strokeWidth: 2)
                      : const Icon(
                        Icons.play_circle_outline,
                        size: 48,
                        color: Colors.black45,
                      ),
            ),
          ),
        ],
      );
    } else {
      child = Stack(
        alignment: Alignment.center,
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: AspectRatio(
              aspectRatio:
                  ctrl.value.aspectRatio == 0 ? 16 / 9 : ctrl.value.aspectRatio,
              child: VideoPlayer(ctrl),
            ),
          ),
          if (!_playing)
            Container(
              color: Colors.black38,
              child: const Icon(
                Icons.play_arrow,
                size: 56,
                color: Colors.white,
              ),
            ),
        ],
      );
    }
    return GestureDetector(
      onTap: _initAndPlay,
      child: SizedBox(width: size.width, height: size.height, child: child),
    );
  }
}
