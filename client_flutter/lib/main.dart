import 'package:flutter/material.dart';
import 'package:client_flutter/theme/app_theme.dart';
import 'package:client_flutter/core/network/health_service.dart';
import 'package:client_flutter/features/auth/login_page.dart';
import 'package:client_flutter/features/home/home_tabs.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  bool _loggedIn = false; // fake auth state
  bool _checkedHealth = false;
  bool _healthOk = false;

  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<void> _init() async {
    // ping backend /health to verify network is ok
    final ok = await HealthService().ping();
    setState(() {
      _healthOk = ok;
      _checkedHealth = true;
    });
  }

  void _onLogin() {
    setState(() => _loggedIn = true);
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Alice Client',
      theme: AppTheme.light(),
      home:
          !_checkedHealth
              ? const _Splash()
              : !_healthOk
              ? _HealthError(onRetry: _init)
              : !_loggedIn
              ? LoginPage(onLogin: _onLogin)
              : const HomeTabs(),
    );
  }
}

class _Splash extends StatelessWidget {
  const _Splash();
  @override
  Widget build(BuildContext context) {
    return const Scaffold(body: Center(child: CircularProgressIndicator()));
  }
}

class _HealthError extends StatelessWidget {
  const _HealthError({required this.onRetry});
  final Future<void> Function() onRetry;
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('服务不可用')),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('无法连接到后端服务 (/health 失败)'),
            const SizedBox(height: 12),
            FilledButton(onPressed: onRetry, child: const Text('重试')),
          ],
        ),
      ),
    );
  }
}
