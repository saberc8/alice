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
