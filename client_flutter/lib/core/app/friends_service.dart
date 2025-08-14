import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/dio_client.dart';

class FriendsService {
  final Dio _dio = DioClient().dio;

  Future<void> sendFriendRequest(String friendEmail) async {
    final res = await _dio.post(
      '/api/v1/app/friends/request',
      data: {'friend_email': friendEmail.trim().toLowerCase()},
    );
    if (res.statusCode != 200) {
      throw Exception('发送请求失败: ${res.statusCode}');
    }
  }

  Future<Map<String, dynamic>> getFriends({
    int page = 1,
    int pageSize = 20,
  }) async {
    final res = await _dio.get(
      '/api/v1/app/friends',
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      if (data is Map<String, dynamic>) return data;
      throw Exception('响应格式错误');
    }
    throw Exception('获取好友列表失败: ${res.statusCode}');
  }

  Future<Map<String, dynamic>> getPendingRequests({
    int page = 1,
    int pageSize = 20,
  }) async {
    final res = await _dio.get(
      '/api/v1/app/friends/requests',
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      if (data is Map<String, dynamic>) return data;
      throw Exception('响应格式错误');
    }
    throw Exception('获取待处理请求失败: ${res.statusCode}');
  }

  Future<void> acceptRequest(int requestId) async {
    final res = await _dio.post(
      '/api/v1/app/friends/requests/$requestId/accept',
    );
    if (res.statusCode != 200) {
      throw Exception('接受失败: ${res.statusCode}');
    }
  }

  Future<void> declineRequest(int requestId) async {
    final res = await _dio.post(
      '/api/v1/app/friends/requests/$requestId/decline',
    );
    if (res.statusCode != 200) {
      throw Exception('拒绝失败: ${res.statusCode}');
    }
  }
}
