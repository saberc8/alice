import 'package:client_flutter/core/app/friends_service.dart';
import 'package:client_flutter/core/network/api_client.dart';
import 'package:flutter/material.dart';

class CreateGroupPage extends StatefulWidget {
  const CreateGroupPage({super.key});
  @override
  State<CreateGroupPage> createState() => _CreateGroupPageState();
}

class _CreateGroupPageState extends State<CreateGroupPage> {
  final _friendsSvc = FriendsService();
  final _api = ApiClient.instance;
  final _nameCtrl = TextEditingController();
  final Set<int> _selected = {};
  bool _loading = true;
  bool _creating = false;
  List<Map<String, dynamic>> _friends = [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final data = await _friendsSvc.getFriends(page: 1, pageSize: 100);
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      _friends = items.cast<Map<String, dynamic>>();
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _create() async {
    if (_creating) return;
    final name = _nameCtrl.text.trim();
    if (name.isEmpty) {
      _show('请输入群名称');
      return;
    }
    if (_selected.length < 2) {
      _show('至少选择2位好友');
      return;
    }
    setState(() => _creating = true);
    try {
      final body = {'name': name, 'member_ids': _selected.toList()};
      final resp = await _api.post(
        '/api/v1/app/chat/groups',
        body: body,
        parser: (d) => d,
      );
      if (!mounted) return;
      Navigator.pop(context, resp);
      // TODO: navigate to group chat page in future
    } catch (e) {
      _show('创建失败: $e');
    } finally {
      if (mounted) setState(() => _creating = false);
    }
  }

  void _show(String msg) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg)));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('发起群聊'),
        actions: [
          TextButton(
            onPressed: _creating ? null : _create,
            child:
                _creating
                    ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                    : const Text('创建'),
          ),
        ],
      ),
      body:
          _loading
              ? const Center(child: CircularProgressIndicator())
              : Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.all(12.0),
                    child: TextField(
                      controller: _nameCtrl,
                      decoration: const InputDecoration(labelText: '群聊名称'),
                    ),
                  ),
                  Expanded(
                    child: ListView.builder(
                      itemCount: _friends.length,
                      itemBuilder: (ctx, i) {
                        final f = _friends[i];
                        final id = (f['id'] as num?)?.toInt() ?? 0;
                        final nick = f['nickname'] ?? f['email'] ?? '用户$id';
                        final avatar = f['avatar'] as String? ?? '';
                        final selected = _selected.contains(id);
                        return ListTile(
                          leading: CircleAvatar(
                            backgroundImage:
                                avatar.isNotEmpty ? NetworkImage(avatar) : null,
                            child:
                                avatar.isEmpty
                                    ? const Icon(Icons.person_outline)
                                    : null,
                          ),
                          title: Text('$nick'),
                          trailing: Checkbox(
                            value: selected,
                            onChanged: (v) {
                              setState(() {
                                if (v == true)
                                  _selected.add(id);
                                else
                                  _selected.remove(id);
                              });
                            },
                          ),
                          onTap: () {
                            setState(() {
                              if (selected)
                                _selected.remove(id);
                              else
                                _selected.add(id);
                            });
                          },
                        );
                      },
                    ),
                  ),
                ],
              ),
    );
  }
}
