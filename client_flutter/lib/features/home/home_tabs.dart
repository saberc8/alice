import 'package:flutter/material.dart';

class HomeTabs extends StatefulWidget {
  const HomeTabs({super.key});

  @override
  State<HomeTabs> createState() => _HomeTabsState();
}

class _HomeTabsState extends State<HomeTabs> {
  int _index = 0;

  final _pages = const [
    _DummyPage(title: '小绿书'),
    _DummyPage(title: '通讯录'),
    _DummyPage(title: '发现'),
    _DummyPage(title: '我'),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_index],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _index,
        type: BottomNavigationBarType.fixed,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.book_outlined),
            label: '小绿书',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.contacts_outlined),
            label: '通讯录',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.explore_outlined),
            label: '发现',
          ),
          BottomNavigationBarItem(icon: Icon(Icons.person_outline), label: '我'),
        ],
        onTap: (i) => setState(() => _index = i),
      ),
    );
  }
}

class _DummyPage extends StatelessWidget {
  const _DummyPage({required this.title});
  final String title;
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(title)),
      body: Center(child: Text('$title 页面')),
    );
  }
}
