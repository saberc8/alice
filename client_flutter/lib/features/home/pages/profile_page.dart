import 'package:flutter/material.dart';
import 'package:client_flutter/core/app/profile_service.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/ui/we_colors.dart';
import 'package:client_flutter/features/profile/edit_profile_page.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key, required this.onLogout});
  final VoidCallback onLogout;
  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: WeAppBar(
        title: '我',
        actions: [
          IconButton(onPressed: _load, icon: const Icon(Icons.refresh)),
          IconButton(
            onPressed: widget.onLogout,
            icon: const Icon(Icons.logout),
          ),
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
                      children: [
                        const WeCell(
                          title: '支付',
                          leading: Icon(
                            Icons.payment_outlined,
                            color: Colors.black87,
                          ),
                        ),
                        const Divider(height: 1, color: WeColors.divider),
                        const WeCell(
                          title: '收藏',
                          leading: Icon(
                            Icons.star_outline,
                            color: Colors.black87,
                          ),
                        ),
                        const Divider(height: 1, color: WeColors.divider),
                        WeCell(
                          title: '编辑资料',
                          leading: const Icon(
                            Icons.edit_outlined,
                            color: Colors.black87,
                          ),
                          onTap: () async {
                            if (_profile == null) return;
                            final updated = await Navigator.of(context).push(
                              MaterialPageRoute(
                                builder:
                                    (_) => EditProfilePage(initial: _profile!),
                              ),
                            );
                            if (updated is Map<String, dynamic>) {
                              setState(() => _profile = updated);
                            }
                          },
                        ),
                        const Divider(height: 1, color: WeColors.divider),
                        const WeCell(
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
