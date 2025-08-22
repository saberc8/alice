import 'package:flutter/material.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/ui/we_colors.dart';
import 'package:client_flutter/features/moments/ui/moment_list_page.dart';

class DiscoverPage extends StatelessWidget {
  const DiscoverPage({super.key});
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const WeAppBar(title: '发现'),
      body: ListView(
        children: [
          const SizedBox(height: 8),
          WeCell(
            leading: const Icon(
              Icons.dynamic_feed_outlined,
              color: Colors.blue,
            ),
            title: '朋友圈',
            subtitle: '查看与发布动态',
            onTap:
                () => Navigator.of(context).push(
                  MaterialPageRoute(builder: (_) => const MomentListPage()),
                ),
          ),
          const Divider(height: 1, color: WeColors.divider),
        ],
      ),
    );
  }
}
