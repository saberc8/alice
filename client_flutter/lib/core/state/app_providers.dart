import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:client_flutter/core/app/friends_service.dart';
import 'package:client_flutter/core/chat/chat_service.dart';
import 'package:client_flutter/core/app/profile_service.dart';
import 'package:client_flutter/core/auth/token_store.dart';

// Token / Auth 状态
final authTokenProvider = StateProvider<String?>(
  (ref) => TokenStore.instance.token,
);

// Service 单例 Provider（可替换为 mock）
final friendsServiceProvider = Provider<FriendsService>(
  (ref) => FriendsService(),
);
final chatServiceProvider = Provider<ChatService>((ref) => ChatService());
final profileServiceProvider = Provider<ProfileService>(
  (ref) => ProfileService(),
);

// 会话列表（演示 Riverpod 使用，真实分页仍可落在页面基类逻辑）
final conversationsProvider =
    FutureProvider.autoDispose<List<Map<String, dynamic>>>((ref) async {
      final chat = ref.watch(chatServiceProvider);
      final data = await chat.getConversations(page: 1, pageSize: 20);
      final items = (data['items'] as List?)?.cast<Map>() ?? [];
      return items.cast<Map<String, dynamic>>();
    });
