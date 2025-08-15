import 'package:flutter/material.dart';
import 'package:client_flutter/theme/app_theme.dart';
// removed health check
import 'package:client_flutter/features/auth/login_page.dart';
import 'package:client_flutter/features/home/home_tabs.dart';
import 'package:client_flutter/core/auth/token_store.dart';
import 'package:client_flutter/features/chat/chat_page.dart';
import 'package:client_flutter/features/contacts/friend_profile_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  bool _loggedIn = false; // based on token presence

  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<void> _init() async {
    await TokenStore.instance.init();
    if (!mounted) return;
    setState(() => _loggedIn = TokenStore.instance.token != null);
  }

  void _onLogin() {
    setState(() => _loggedIn = true);
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Alice Client',
      theme: AppTheme.light(),
      debugShowCheckedModeBanner: false,
      onGenerateRoute: (settings) {
        if (settings.name == '/friend_profile') {
          final user = settings.arguments as Map<String, dynamic>;
          return MaterialPageRoute(
            builder: (_) => FriendProfilePage(user: user),
            settings: settings,
          );
        }
        if (settings.name == '/chat') {
          final user = settings.arguments as Map<String, dynamic>;
          return MaterialPageRoute(
            builder: (_) => ChatPage(peer: user),
            settings: settings,
          );
        }
        return null;
      },
      home:
          !_loggedIn
              ? LoginPage(onLogin: _onLogin)
              : HomeTabs(
                onLogout: () async {
                  await TokenStore.instance.clear();
                  if (!mounted) return;
                  setState(() => _loggedIn = false);
                },
              ),
    );
  }
}
