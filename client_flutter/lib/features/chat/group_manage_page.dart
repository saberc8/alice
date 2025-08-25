import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:client_flutter/core/chat/chat_service.dart';

class GroupManagePage extends StatefulWidget {
  const GroupManagePage({super.key, required this.group});
  final Map<String, dynamic> group; // {group_id / id, name, avatar}
  @override
  State<GroupManagePage> createState() => _GroupManagePageState();
}

class _GroupManagePageState extends State<GroupManagePage> {
  final _svc = ChatService();
  bool _saving = false;
  late TextEditingController _nameCtrl;
  String? _avatar; // local preview (url or path)
  List<Map<String, dynamic>> _members = [];
  bool _loadingMembers = false;

  int get _groupId {
    final v = widget.group['group_id'] ?? widget.group['id'];
    if (v is int) return v;
    if (v is double) return v.toInt();
    if (v is String) return int.tryParse(v) ?? 0;
    return 0;
  }

  @override
  void initState() {
    super.initState();
    _nameCtrl = TextEditingController(
      text: widget.group['name'] ?? widget.group['nickname'] ?? '',
    );
    _avatar = widget.group['avatar'] as String?;
    _loadMembers();
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    super.dispose();
  }

  Future<void> _loadMembers() async {
    setState(() => _loadingMembers = true);
    try {
      final list = await _svc.listGroupMembers(_groupId);
      list.sort((a, b) => (a['id'] ?? 0).compareTo(b['id'] ?? 0));
      setState(() => _members = list);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('加载成员失败: $e')));
      }
    } finally {
      if (mounted) setState(() => _loadingMembers = false);
    }
  }

  Future<void> _pickAvatar() async {
    final picker = ImagePicker();
    final img = await picker.pickImage(
      source: ImageSource.gallery,
      imageQuality: 85,
    );
    if (img == null) return;
    setState(() => _saving = true);
    try {
      final pathOrUrl = await _svc.uploadGroupAvatar(_groupId, img.path);
      if (pathOrUrl != null) {
        setState(() => _avatar = pathOrUrl);
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _save() async {
    final name = _nameCtrl.text.trim();
    if (name.isEmpty) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('群名称不能为空')));
      return;
    }
    setState(() => _saving = true);
    try {
      await _svc.updateGroup(groupId: _groupId, name: name, avatar: _avatar);
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('保存成功')));
        Navigator.pop(context, {
          'updated': true,
          'name': name,
          'avatar': _avatar,
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('保存失败: $e')));
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _addMembersDialog() async {
    // 简化：输入用户ID列表，用逗号分隔
    final ctrl = TextEditingController();
    final ok = await showDialog<bool>(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text('添加成员'),
            content: TextField(
              controller: ctrl,
              decoration: const InputDecoration(hintText: '输入用户ID，例: 2,3,5'),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text('取消'),
              ),
              TextButton(
                onPressed: () => Navigator.pop(ctx, true),
                child: const Text('确定'),
              ),
            ],
          ),
    );
    if (ok != true) return;
    final text = ctrl.text.trim();
    if (text.isEmpty) return;
    final ids =
        text
            .split(',')
            .map((e) => int.tryParse(e.trim()))
            .whereType<int>()
            .toList();
    if (ids.isEmpty) return;
    try {
      await _svc.addGroupMembers(_groupId, ids);
      await _loadMembers();
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('添加成功')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('添加失败: $e')));
      }
    }
  }

  Future<void> _removeMember(int userId) async {
    final ok = await showDialog<bool>(
      context: context,
      builder:
          (ctx) => AlertDialog(
            title: const Text('移除成员'),
            content: Text('确定移除用户 $userId ?'),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text('取消'),
              ),
              TextButton(
                onPressed: () => Navigator.pop(ctx, true),
                child: const Text('移除'),
              ),
            ],
          ),
    );
    if (ok != true) return;
    try {
      await _svc.removeGroupMember(_groupId, userId);
      await _loadMembers();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('移除失败: $e')));
      }
    }
  }

  Widget _buildMembersSection() {
    if (_loadingMembers) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(16),
          child: CircularProgressIndicator(),
        ),
      );
    }
    if (_members.isEmpty) {
      return const Padding(padding: EdgeInsets.all(16), child: Text('暂无成员'));
    }
    return ListView.separated(
      itemCount: _members.length,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      separatorBuilder: (_, __) => const Divider(height: 1),
      itemBuilder: (ctx, i) {
        final m = _members[i];
        final avatar = m['avatar'] as String?;
        final name = m['nickname'] ?? m['name'] ?? '用户${m['id']}';
        final id = m['id'] ?? 0;
        return ListTile(
          leading: CircleAvatar(
            backgroundImage:
                (avatar != null && avatar.startsWith('http'))
                    ? NetworkImage(avatar)
                    : null,
            child:
                (avatar == null || avatar.isEmpty)
                    ? Text(name.isNotEmpty ? name[0] : '?')
                    : null,
          ),
          title: Text(name),
          subtitle: Text('ID: $id'),
          trailing: IconButton(
            icon: const Icon(
              Icons.remove_circle_outline,
              color: Colors.redAccent,
            ),
            onPressed: () => _removeMember(id),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('群管理'),
        actions: [
          TextButton(
            onPressed: _saving ? null : _save,
            child:
                _saving
                    ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                    : const Text('保存'),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _addMembersDialog,
        child: const Icon(Icons.person_add_alt_1),
      ),
      body: AbsorbPointer(
        absorbing: _saving,
        child: ListView(
          padding: const EdgeInsets.symmetric(vertical: 12),
          children: [
            Center(
              child: Stack(
                children: [
                  CircleAvatar(
                    radius: 42,
                    backgroundImage:
                        (_avatar != null &&
                                _avatar!.isNotEmpty &&
                                (_avatar!.startsWith('http://') ||
                                    _avatar!.startsWith('https://')))
                            ? NetworkImage(_avatar!)
                            : null,
                    child:
                        (_avatar == null || _avatar!.isEmpty)
                            ? const Icon(Icons.group_outlined, size: 42)
                            : null,
                  ),
                  Positioned(
                    bottom: 0,
                    right: 0,
                    child: InkWell(
                      onTap: _pickAvatar,
                      child: Container(
                        padding: const EdgeInsets.all(6),
                        decoration: BoxDecoration(
                          color: Colors.black54,
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: const Icon(
                          Icons.edit,
                          color: Colors.white,
                          size: 16,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 12),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: TextField(
                controller: _nameCtrl,
                maxLength: 30,
                decoration: const InputDecoration(labelText: '群名称'),
              ),
            ),
            const Divider(height: 32),
            const ListTile(
              leading: Icon(Icons.people_outline),
              title: Text('成员列表'),
            ),
            _buildMembersSection(),
          ],
        ),
      ),
    );
  }
}
