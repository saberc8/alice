import 'package:flutter/material.dart';
import 'we_colors.dart';

class WeBadge extends StatelessWidget {
  const WeBadge({super.key, required this.child, this.count});
  final Widget child;
  final int? count;

  @override
  Widget build(BuildContext context) {
    if (count == null || count! <= 0) return child;
    final text = count! > 99 ? '99+' : '$count';
    return Stack(
      clipBehavior: Clip.none,
      children: [
        child,
        Positioned(
          right: -2,
          top: -2,
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: WeColors.badge,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.white, width: 1),
            ),
            child: Text(
              text,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 10,
                height: 1.0,
              ),
            ),
          ),
        ),
      ],
    );
  }
}
