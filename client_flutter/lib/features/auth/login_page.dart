import 'package:flutter/material.dart';
import 'package:client_flutter/theme/app_theme.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key, required this.onLogin});

  final void Function() onLogin;

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _userCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  bool _loading = false;

  @override
  void dispose() {
    _userCtrl.dispose();
    _pwdCtrl.dispose();
    super.dispose();
  }

  Future<void> _fakeLogin() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _loading = true);
    await Future.delayed(const Duration(milliseconds: 800));
    if (!mounted) return;
    widget.onLogin();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('登录')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 24),
              TextFormField(
                controller: _userCtrl,
                decoration: const InputDecoration(labelText: '用户名'),
                validator: (v) => (v == null || v.isEmpty) ? '请输入用户名' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _pwdCtrl,
                decoration: const InputDecoration(labelText: '密码'),
                obscureText: true,
                validator: (v) => (v == null || v.isEmpty) ? '请输入密码' : null,
              ),
              const SizedBox(height: 24),
              FilledButton(
                onPressed: _loading ? null : _fakeLogin,
                style: FilledButton.styleFrom(
                  backgroundColor: AppTheme.primary,
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child:
                    _loading
                        ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                        : const Text(
                          '登录',
                          style: TextStyle(fontSize: 16, color: Colors.white),
                        ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
