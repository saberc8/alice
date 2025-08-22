import 'package:client_flutter/core/network/api_client.dart';
import 'package:client_flutter/core/auth/token_store.dart';

class AuthService {
  final _api = ApiClient.instance;

  Future<void> register({
    required String email,
    required String password,
    String? nickname,
  }) async {
    final data = await _api.post<Map<String, dynamic>>(
      '/api/v1/app/register',
      body: {
        'email': email.trim().toLowerCase(),
        'password': password,
        if (nickname != null && nickname.isNotEmpty) 'nickname': nickname,
      },
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
    final token = data['token'] as String?;
    if (token == null || token.isEmpty) {
      throw ApiException('注册成功但未返回 token');
    }
    await TokenStore.instance.setToken(token);
  }

  Future<void> login({required String email, required String password}) async {
    final data = await _api.post<Map<String, dynamic>>(
      '/api/v1/app/login',
      body: {'email': email.trim().toLowerCase(), 'password': password},
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
    final token = data['token'] as String?;
    if (token == null || token.isEmpty) {
      throw ApiException('登录成功但未返回 token');
    }
    await TokenStore.instance.setToken(token);
  }
}
