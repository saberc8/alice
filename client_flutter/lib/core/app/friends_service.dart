import 'package:client_flutter/core/network/api_client.dart';

class FriendsService {
  final _api = ApiClient.instance;

  Future<void> sendFriendRequest(String friendEmail) async {
    await _api.post(
      '/api/v1/app/friends/request',
      body: {'friend_email': friendEmail.trim().toLowerCase()},
    );
  }

  Future<Map<String, dynamic>> getFriends({
    int page = 1,
    int pageSize = 20,
  }) async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/friends',
      query: {'page': page, 'page_size': pageSize},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  Future<Map<String, dynamic>> getPendingRequests({
    int page = 1,
    int pageSize = 20,
  }) async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/friends/requests',
      query: {'page': page, 'page_size': pageSize},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  Future<void> acceptRequest(int requestId) async {
    await _api.post('/api/v1/app/friends/requests/$requestId/accept');
  }

  Future<void> declineRequest(int requestId) async {
    await _api.post('/api/v1/app/friends/requests/$requestId/decline');
  }
}
