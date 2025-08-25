import 'package:flutter/material.dart';

/// 轻量级 Emoji 选择组件
/// 1. 纯 Unicode，不涉及版权图片，可直接存库或传输
/// 2. 支持简单分类 Tab，可扩展为网络/数据库加载
/// 3. 外部通过 [onSelected] 获取选择结果
class EmojiPicker extends StatefulWidget {
  const EmojiPicker({
    super.key,
    this.onSelected,
    this.backgroundColor,
    this.crossAxisCount = 8,
    this.emojiSize = 28,
    this.padding = const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
  });

  final ValueChanged<String>? onSelected;
  final Color? backgroundColor;
  final int crossAxisCount;
  final double emojiSize;
  final EdgeInsets padding;

  @override
  State<EmojiPicker> createState() => _EmojiPickerState();
}

class _EmojiPickerState extends State<EmojiPicker>
    with SingleTickerProviderStateMixin {
  late final TabController _tabController;

  // 基础分类数据，可后续替换为后端下发
  final List<(String label, List<String> emojis)> _groups = [
    (
      '常用',
      [
        '😀',
        '😁',
        '😂',
        '🤣',
        '😃',
        '😄',
        '😅',
        '🥹',
        '😊',
        '😉',
        '😍',
        '😘',
        '😗',
        '😙',
        '😚',
        '😋',
        '😛',
        '😜',
        '🤪',
        '🤨',
        '🫠',
        '🤔',
        '😶',
        '🫥',
        '🙂',
        '🙃',
        '😇',
        '🥰',
        '😭',
        '😤',
        '😡',
        '👍',
        '👎',
        '👌',
        '🙏',
        '👏',
        '🤝',
        '👀',
        '💪',
        '🔥',
        '⭐',
        '🌟',
        '💯',
        '✨',
        '❤️',
        '🧡',
        '💛',
        '💚',
        '💙',
        '💜',
        '🖤',
        '🤍',
        '🤎',
        '💔',
        '❣️',
        '💕',
        '💞',
        '💓',
        '💗',
      ],
    ),
    (
      '人物',
      [
        '👋',
        '🤚',
        '🖐',
        '✋',
        '🖖',
        '👌',
        '🤌',
        '🤏',
        '✌️',
        '🤞',
        '🤟',
        '🤘',
        '🤙',
        '👈',
        '👉',
        '👆',
        '🖕',
        '👇',
        '☝️',
        '👍',
        '👎',
        '✊',
        '👊',
        '🤛',
        '🤜',
        '👏',
        '🙌',
        '👐',
        '🤲',
        '🙏',
        '✍️',
        '💅',
        '🤳',
        '💪',
        '🦾',
        '🦵',
        '🦿',
        '🦶',
        '👣',
        '👂',
        '🦻',
        '👃',
        '👀',
        '👁',
        '🧠',
        '🫀',
        '🫁',
        '🦷',
        '🦴',
        '👅',
        '👄',
        '🫦',
        '👶',
        '🧒',
        '👦',
        '👧',
        '🧑',
        '👱',
        '👨',
        '👩',
        '🧔',
        '🧔‍♂️',
        '🧔‍♀️',
        '👨‍🦰',
        '👩‍🦰',
        '👨‍🦱',
        '👩‍🦱',
      ],
    ),
    (
      '自然',
      [
        '🌞',
        '🌝',
        '🌛',
        '⭐',
        '🌟',
        '✨',
        '⚡',
        '🔥',
        '🌈',
        '☁️',
        '⛅',
        '🌧',
        '⛈',
        '🌩',
        '❄️',
        '💧',
        '💦',
        '🌊',
        '🪨',
        '🌵',
        '🌲',
        '🌳',
        '🌴',
        '🌱',
        '🌿',
        '☘️',
        '🍀',
        '🎍',
        '🌷',
        '🌹',
        '🥀',
        '🌺',
        '🌸',
        '🌼',
        '🌻',
        '🍄',
        '🪺',
        '🦋',
      ],
    ),
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _groups.length, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final bg = widget.backgroundColor ?? Theme.of(context).colorScheme.surface;
    return Material(
      color: bg,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
            height: 40,
            child: TabBar(
              controller: _tabController,
              isScrollable: true,
              labelPadding: const EdgeInsets.symmetric(horizontal: 12),
              tabAlignment: TabAlignment.start,
              indicatorSize: TabBarIndicatorSize.label,
              tabs: [for (final g in _groups) Tab(text: g.$1, height: 36)],
            ),
          ),
          SizedBox(
            height: 250, // 固定高度，避免键盘弹出布局跳动
            child: TabBarView(
              controller: _tabController,
              children: [
                for (final g in _groups)
                  GridView.builder(
                    padding: widget.padding,
                    gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: widget.crossAxisCount,
                      mainAxisSpacing: 4,
                      crossAxisSpacing: 4,
                    ),
                    itemCount: g.$2.length,
                    itemBuilder: (ctx, i) {
                      final e = g.$2[i];
                      return InkWell(
                        borderRadius: BorderRadius.circular(8),
                        onTap: () => widget.onSelected?.call(e),
                        child: Center(
                          child: Text(
                            e,
                            style: TextStyle(fontSize: widget.emojiSize),
                          ),
                        ),
                      );
                    },
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
