import 'dart:io';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:client_flutter/core/app/profile_service.dart';

class EditProfilePage extends StatefulWidget {
  const EditProfilePage({super.key, required this.initial});
  final Map<String, dynamic> initial;

  @override
  State<EditProfilePage> createState() => _EditProfilePageState();
}

class _EditProfilePageState extends State<EditProfilePage> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nickname;
  late TextEditingController _bio;
  String? _gender; // male female other
  bool _submitting = false;
  File? _avatarFile;
  Map<String, dynamic>? _current;

  final _svc = ProfileService();

  @override
  void initState() {
    super.initState();
    _current = Map<String, dynamic>.from(widget.initial);
    _nickname = TextEditingController(text: widget.initial['nickname'] ?? '');
    _bio = TextEditingController(text: widget.initial['bio'] ?? '');
    _gender = widget.initial['gender'] as String?;
  }

  @override
  void dispose() {
    _nickname.dispose();
    _bio.dispose();
    super.dispose();
  }

  Future<void> _pickAvatar() async {
    final picker = ImagePicker();
    final x = await picker.pickImage(
      source: ImageSource.gallery,
      imageQuality: 85,
    );
    if (x == null) return;
    setState(() => _avatarFile = File(x.path));
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _submitting = true);
    try {
      Map<String, dynamic>? updated;
      if (_avatarFile != null) {
        updated = await _svc.uploadAvatar(_avatarFile!);
      }
      updated = await _svc.updateProfile(
        nickname: _nickname.text.trim().isEmpty ? null : _nickname.text.trim(),
        gender: _gender,
        bio: _bio.text.trim().isEmpty ? null : _bio.text.trim(),
      );
      if (!mounted) return;
      Navigator.of(context).pop(updated);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('保存失败: $e')));
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  void _changePassword() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => const _ChangePasswordSheet(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final avatarUrl = _current?['avatar'] as String?;
    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑资料'),
        actions: [
          TextButton(
            onPressed: _submitting ? null : _save,
            child:
                _submitting
                    ? const SizedBox(
                      height: 18,
                      width: 18,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                    : const Text('保存'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Center(
              child: Stack(
                children: [
                  CircleAvatar(
                    radius: 46,
                    backgroundImage:
                        _avatarFile != null
                            ? FileImage(_avatarFile!)
                            : (avatarUrl != null && avatarUrl.isNotEmpty)
                            ? NetworkImage(avatarUrl)
                            : null,
                    child:
                        (_avatarFile == null &&
                                (avatarUrl == null || avatarUrl.isEmpty))
                            ? const Icon(Icons.person, size: 46)
                            : null,
                  ),
                  Positioned(
                    bottom: 0,
                    right: 0,
                    child: Material(
                      color: Colors.black54,
                      shape: const CircleBorder(),
                      child: IconButton(
                        icon: const Icon(
                          Icons.edit,
                          size: 20,
                          color: Colors.white,
                        ),
                        onPressed: _pickAvatar,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),
            TextFormField(
              controller: _nickname,
              decoration: const InputDecoration(labelText: '昵称'),
              maxLength: 30,
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              value: _gender,
              decoration: const InputDecoration(labelText: '性别'),
              items: const [
                DropdownMenuItem(value: 'male', child: Text('男')),
                DropdownMenuItem(value: 'female', child: Text('女')),
                DropdownMenuItem(value: 'other', child: Text('其他')),
              ],
              onChanged: (v) => setState(() => _gender = v),
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _bio,
              decoration: const InputDecoration(labelText: '个性签名'),
              maxLines: 3,
              maxLength: 160,
            ),
            const SizedBox(height: 24),
            ListTile(
              leading: const Icon(Icons.lock_outline),
              title: const Text('修改密码'),
              trailing: const Icon(Icons.chevron_right),
              onTap: _changePassword,
            ),
          ],
        ),
      ),
    );
  }
}

class _ChangePasswordSheet extends StatefulWidget {
  const _ChangePasswordSheet();
  @override
  State<_ChangePasswordSheet> createState() => _ChangePasswordSheetState();
}

class _ChangePasswordSheetState extends State<_ChangePasswordSheet> {
  final _formKey = GlobalKey<FormState>();
  final _old = TextEditingController();
  final _new = TextEditingController();
  final _new2 = TextEditingController();
  bool _submitting = false;
  final _svc = ProfileService();

  @override
  void dispose() {
    _old.dispose();
    _new.dispose();
    _new2.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_new.text != _new2.text) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('两次新密码不一致')));
      return;
    }
    setState(() => _submitting = true);
    try {
      await _svc.changePassword(oldPassword: _old.text, newPassword: _new.text);
      if (!mounted) return;
      Navigator.of(context).pop();
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('密码已修改')));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('修改失败: $e')));
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final viewInsets = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.only(bottom: viewInsets),
      child: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextFormField(
                  controller: _old,
                  decoration: const InputDecoration(labelText: '旧密码'),
                  obscureText: true,
                  validator: (v) => (v == null || v.isEmpty) ? '请输入旧密码' : null,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _new,
                  decoration: const InputDecoration(labelText: '新密码'),
                  obscureText: true,
                  validator: (v) => (v == null || v.length < 6) ? '至少6位' : null,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _new2,
                  decoration: const InputDecoration(labelText: '重复新密码'),
                  obscureText: true,
                  validator: (v) => (v == null || v.length < 6) ? '至少6位' : null,
                ),
                const SizedBox(height: 20),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton(
                    onPressed: _submitting ? null : _submit,
                    child:
                        _submitting
                            ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                            : const Text('确认修改'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
