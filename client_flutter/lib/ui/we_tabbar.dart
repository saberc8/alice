import 'package:flutter/material.dart';
import 'we_colors.dart';
import 'we_badge.dart';

class WeTabItem {
  const WeTabItem({
    required this.icon,
    required this.iconActive,
    required this.label,
    this.badgeCount,
  });

  final IconData icon;
  final IconData iconActive;
  final String label;
  final int? badgeCount;
}

class WeTabBar extends StatelessWidget {
  const WeTabBar({
    super.key,
    required this.items,
    required this.currentIndex,
    required this.onTap,
  });

  final List<WeTabItem> items;
  final int currentIndex;
  final ValueChanged<int> onTap;

  @override
  Widget build(BuildContext context) {
    // 处理底部安全区 (iPhone 刘海/底部指示条 等)，避免出现额外空白或颜色不一致的小 bar。
    final bottomInset = MediaQuery.paddingOf(context).bottom;
    final effectiveHeight = 56.0 + bottomInset; // 让点击区域延伸到底部，同时视觉高度保持 56
    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        border: Border(top: BorderSide(color: Color(0xFFE5E5E5), width: 0.5)),
      ),
      height: effectiveHeight,
      padding: EdgeInsets.only(bottom: bottomInset),
      child: Row(
        children: [
          for (int i = 0; i < items.length; i++) _buildItem(context, i),
        ],
      ),
    );
  }

  Widget _buildItem(BuildContext context, int index) {
    final it = items[index];
    final selected = index == currentIndex;
    final color = selected ? WeColors.brand : WeColors.textSecondary;
    return Expanded(
      child: InkWell(
        onTap: () => onTap(index),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            WeBadge(
              count: it.badgeCount,
              child: Icon(selected ? it.iconActive : it.icon, color: color),
            ),
            const SizedBox(height: 2),
            Text(it.label, style: TextStyle(fontSize: 11, color: color)),
          ],
        ),
      ),
    );
  }
}
