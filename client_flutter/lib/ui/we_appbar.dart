import 'package:flutter/material.dart';
import 'we_colors.dart';

/// A minimalist AppBar similar to WeChat
class WeAppBar extends StatelessWidget implements PreferredSizeWidget {
  const WeAppBar({super.key, required this.title, this.actions});
  final String title;
  final List<Widget>? actions;

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    final bg = Theme.of(context).scaffoldBackgroundColor;
    return AppBar(
      centerTitle: false,
      titleSpacing: 16,
      backgroundColor: bg,
      elevation: 0.4,
      surfaceTintColor: Colors.transparent,
      shadowColor: Colors.black12,
      title: Text(
        title,
        style: const TextStyle(
          color: WeColors.textPrimary,
          fontWeight: FontWeight.w600,
          fontSize: 18,
        ),
      ),
      actions: actions,
      iconTheme: const IconThemeData(color: WeColors.textPrimary),
    );
  }
}
