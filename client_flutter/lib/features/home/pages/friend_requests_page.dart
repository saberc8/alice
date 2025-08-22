import 'package:flutter/material.dart';
import 'package:client_flutter/core/app/friends_service.dart';
import 'package:client_flutter/ui/we_appbar.dart';
import 'package:client_flutter/ui/we_cell.dart';
import 'package:client_flutter/core/util/base_list_page_state.dart';

class FriendRequestsPage extends StatefulWidget {
  const FriendRequestsPage({super.key});
  @override
  State<FriendRequestsPage> createState() => _FriendRequestsPageState();
}

class _FriendRequestsPageState
    extends BaseListPageState<_FriendRequestItem, FriendRequestsPage> {
  final _svc = FriendsService();

  @override
  Future<List<_FriendRequestItem>> fetch({
    required int page,
    required int pageSize,
  }) async {
    final data = await _svc.getPendingRequests(page: page, pageSize: pageSize);
    final ids = (data['request_ids'] as List?)?.cast<int>() ?? [];
    final requesterIds = (data['requester_ids'] as List?)?.cast<int>() ?? [];
    final list = <_FriendRequestItem>[];
    for (var i = 0; i < ids.length; i++) {
      list.add(
        _FriendRequestItem(
          requestId: ids[i],
          requesterId: i < requesterIds.length ? requesterIds[i] : null,
        ),
      );
    }
    return list;
  }

  Future<void> _accept(_FriendRequestItem item) async {
    await _svc.acceptRequest(item.requestId);
    reload();
  }

  Future<void> _decline(_FriendRequestItem item) async {
    await _svc.declineRequest(item.requestId);
    reload();
  }

  @override
  Widget buildItem(BuildContext context, _FriendRequestItem item, int index) {
    return WeCell(
      leading: const CircleAvatar(child: Icon(Icons.person_outline)),
      title: '请求者 ID: ${item.requesterId ?? '-'}',
      subtitle: '请求 ID: ${item.requestId}',
      trailing: Wrap(
        spacing: 8,
        children: [
          TextButton(onPressed: () => _decline(item), child: const Text('拒绝')),
          FilledButton(onPressed: () => _accept(item), child: const Text('接受')),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: WeAppBar(title: '待处理好友请求'),
      body: _FriendRequestsBody(),
    );
  }
}

// 将 Body 分离以便使用父类的 build(context)
class _FriendRequestsBody extends StatelessWidget {
  const _FriendRequestsBody();
  @override
  Widget build(BuildContext context) {
    final state = context.findAncestorStateOfType<_FriendRequestsPageState>();
    if (state == null) return const SizedBox.shrink();
    return state.build(context);
  }
}

class _FriendRequestItem {
  final int requestId;
  final int? requesterId;
  _FriendRequestItem({required this.requestId, this.requesterId});
}
