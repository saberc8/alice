import 'dart:async';
import 'dart:convert';

import 'package:client_flutter/core/network/dio_client.dart';
import 'package:client_flutter/core/network/api_client.dart';
import 'package:client_flutter/core/auth/token_store.dart';
import 'package:dio/dio.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

/// Lightweight chat service for p2p chat via backend described in docs/ws.md
class ChatService {
  ChatService._();
  static final ChatService _singleton = ChatService._();
  factory ChatService() => _singleton;

  final Dio _dio = DioClient().dio;
  final _api = ApiClient.instance;

  int? _selfId; // 当前用户 id
  final _bus = StreamController<Map<String, dynamic>>.broadcast();
  Stream<Map<String, dynamic>> get messageStream => _bus.stream;

  /// Open a websocket connection with Bearer token via query param.
  /// Returns a [Stream] of decoded message maps and a [sink] to send.
  (
    Stream<Map<String, dynamic>> stream,
    StreamSink<Map<String, dynamic>> sink,
    Future<void> Function() close,
  )
  connect() {
    // baseUrl may have path; ensure we build ws/wss url correctly
    final httpBase = Uri.parse(_dio.options.baseUrl);
    final wsScheme = httpBase.scheme == 'https' ? 'wss' : 'ws';

    // Read token directly from TokenStore (Dio BaseOptions.headers won't have it)
    final token = TokenStore.instance.token;
    if (token == null || token.isEmpty) {
      throw Exception('未登录或 token 缺失，无法建立聊天连接');
    }

    final uri = httpBase.replace(
      scheme: wsScheme,
      path: '/api/v1/app/chat/ws',
      queryParameters: {'token': token},
    );

    final channel = WebSocketChannel.connect(uri);

    // Wrap as json stream/sink
    final controller = StreamController<Map<String, dynamic>>.broadcast();
    final sinkController = StreamController<Map<String, dynamic>>();

    // incoming
    channel.stream.listen(
      (event) async {
        try {
          final data = event is String ? jsonDecode(event) : event;
          if (data is Map<String, dynamic>) {
            controller.add(data);
            _bus.add(data);
            if (_selfId == null) _ensureProfile();
          }
        } catch (_) {}
      },
      onError: controller.addError,
      onDone: controller.close,
    );

    // outgoing
    sinkController.stream.listen((msg) {
      channel.sink.add(jsonEncode(msg));
    });

    Future<void> close() async {
      await channel.sink.close();
      await controller.close();
      await sinkController.close();
    }

    _ensureProfile(); // 异步获取当前用户 id
    return (controller.stream, sinkController.sink, close);
  }

  /// Send one message via REST is not required; we send via WS by adding to sink.
  /// Query history
  Future<Map<String, dynamic>> getHistory({
    required int peerId,
    int page = 1,
    int pageSize = 20,
  }) async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/chat/history/$peerId',
      query: {'page': page, 'page_size': pageSize},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  Future<Map<String, dynamic>> getGroupHistory({
    required int groupId,
    int page = 1,
    int pageSize = 20,
  }) async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/chat/groups/$groupId/messages',
      query: {'page': page, 'page_size': pageSize},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  /// 获取最近会话列表
  Future<Map<String, dynamic>> getConversations({
    int page = 1,
    int pageSize = 20,
  }) async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/chat/conversations',
      query: {'page': page, 'page_size': pageSize},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  Future<void> markRead({required int peerId, required int beforeId}) async {
    if (beforeId <= 0) return;
    await _api.post(
      '/api/v1/app/chat/read',
      body: {'peer_id': peerId, 'before_id': beforeId},
    );
  }

  Future<void> markGroupRead({
    required int groupId,
    required int beforeMsgId,
  }) async {
    if (beforeMsgId <= 0) return;
    await _api.post(
      '/api/v1/app/chat/groups/read',
      body: {'group_id': groupId, 'before_msg_id': beforeMsgId},
    );
  }

  Future<Map<String, dynamic>?> updateGroup({
    required int groupId,
    String? name,
    String? avatar,
  }) async {
    final body = <String, dynamic>{};
    if (name != null) body['name'] = name;
    if (avatar != null) body['avatar'] = avatar;
    final resp = await _api.put<Map<String, dynamic>>(
      '/api/v1/app/chat/groups/$groupId',
      body: body,
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
    return resp;
  }

  Future<String?> uploadGroupAvatar(int groupId, String path) async {
    try {
      final fileName = path.split('/').last;
      final formData = FormData.fromMap({
        'file': await MultipartFile.fromFile(path, filename: fileName),
      });
      final resp = await _dio.post(
        '/api/v1/app/chat/groups/$groupId/avatar',
        data: formData,
        options: Options(headers: {'Content-Type': 'multipart/form-data'}),
      );
      final data = resp.data;
      if (data is Map && data['data'] is Map) {
        final inner = data['data'] as Map;
        final pathOrUrl = inner['url'] ?? inner['path'];
        if (pathOrUrl is String) return pathOrUrl;
      }
    } catch (_) {}
    return null;
  }

  // ===== Group Members APIs =====
  Future<List<Map<String, dynamic>>> listGroupMembers(int groupId) async {
    final resp = await _api.get<Map<String, dynamic>>(
      '/api/v1/app/chat/groups/$groupId/members',
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
    final data = resp['members'];
    if (data is List) {
      return data
          .map((e) => (e is Map<String, dynamic>) ? e : <String, dynamic>{})
          .toList();
    }
    return [];
  }

  Future<void> addGroupMembers(int groupId, List<int> userIds) async {
    if (userIds.isEmpty) return;
    await _api.post(
      '/api/v1/app/chat/groups/$groupId/members/add',
      body: {'user_ids': userIds},
    );
  }

  Future<void> removeGroupMember(int groupId, int userId) async {
    await _api.post(
      '/api/v1/app/chat/groups/$groupId/members/remove',
      body: {'user_id': userId},
    );
  }

  Future<int?> selfId() async {
    if (_selfId != null) return _selfId;
    await _ensureProfile();
    return _selfId;
  }

  /// 上传图片到对象存储后返回 url，供发送图片消息使用
  /// bucket 可根据后端策略自定义，这里采用 `chat`，若不存在需提前创建
  Future<String?> uploadImage(String path) async {
    try {
      final fileName = path.split('/').last;
      final formData = FormData.fromMap({
        'file': await MultipartFile.fromFile(path, filename: fileName),
      });
      final resp = await _dio.post(
        '/api/v1/app/chat/images',
        data: formData,
        options: Options(headers: {'Content-Type': 'multipart/form-data'}),
      );
      final data = resp.data;
      if (data is Map && data['data'] is Map) {
        final inner = data['data'] as Map;
        final url = inner['url'] ?? inner['path'];
        if (url is String) return url;
      }
    } catch (_) {}
    return null;
  }

  /// 上传视频
  Future<String?> uploadVideo(String path) async {
    try {
      final fileName = path.split('/').last;
      final formData = FormData.fromMap({
        'file': await MultipartFile.fromFile(path, filename: fileName),
      });
      final resp = await _dio.post(
        '/api/v1/app/chat/videos',
        data: formData,
        options: Options(headers: {'Content-Type': 'multipart/form-data'}),
      );
      final data = resp.data;
      if (data is Map && data['data'] is Map) {
        final inner = data['data'] as Map;
        final url = inner['url'] ?? inner['path'];
        if (url is String) return url;
      }
    } catch (_) {}
    return null;
  }

  Future<void> _ensureProfile() async {
    if (_selfId != null) return;
    try {
      final data = await _api.get<Map<String, dynamic>>(
        '/api/v1/app/profile',
        parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
      );
      final id = data['id'];
      if (id is int)
        _selfId = id;
      else if (id is double)
        _selfId = id.toInt();
    } catch (_) {}
  }
}
