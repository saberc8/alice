import 'package:flutter/material.dart';
import 'package:client_flutter/ui/we_tabbar.dart';
import 'pages/conversations_page.dart';
import 'pages/contacts_page.dart';
import 'pages/discover_page.dart';
import 'pages/profile_page.dart';

class HomeTabs extends StatefulWidget {
  const HomeTabs({super.key, required this.onLogout});

  final VoidCallback onLogout;

  @override
  State<HomeTabs> createState() => _HomeTabsState();
}

class _HomeTabsState extends State<HomeTabs> {
  int _index = 0;
  late final List<Widget> _pages;

  @override
  void initState() {
    super.initState();
    // 缓存页面，避免 tab 切换时重复创建导致状态丢失或重建。
    _pages = [
      const ConversationsPage(),
      const ContactsPage(),
      const DiscoverPage(),
      ProfilePage(onLogout: widget.onLogout),
    ];
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(index: _index, children: _pages),
      bottomNavigationBar: WeTabBar(
        items: const [
          WeTabItem(
            icon: Icons.chat_bubble_outline,
            iconActive: Icons.chat_bubble,
            label: '小绿书',
          ),
          WeTabItem(
            icon: Icons.contacts_outlined,
            iconActive: Icons.contacts,
            label: '通讯录',
          ),
          WeTabItem(
            icon: Icons.explore_outlined,
            iconActive: Icons.explore,
            label: '发现',
          ),
          WeTabItem(
            icon: Icons.person_outline,
            iconActive: Icons.person,
            label: '我',
          ),
        ],
        currentIndex: _index,
        onTap: (i) => setState(() => _index = i),
      ),
    );
  }
}
