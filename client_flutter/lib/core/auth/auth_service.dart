import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/dio_client.dart';
import 'package:client_flutter/core/auth/token_store.dart';

class AuthService {
  final Dio _dio = DioClient().dio;

  Future<void> register({
    required String email,
    required String password,
    String? nickname,
  }) async {
    final res = await _dio.post(
      '/api/v1/app/register',
      data: {
        'email': email.trim().toLowerCase(),
        'password': password,
        if (nickname != null && nickname.isNotEmpty) 'nickname': nickname,
      },
    );
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      final token = data is Map ? data['token'] as String? : null;
      if (token != null && token.isNotEmpty) {
        await TokenStore.instance.setToken(token);
        return;
      }
      throw Exception('注册成功但未返回 token');
    }
    throw Exception('注册失败: ${res.statusCode}');
  }

  Future<void> login({required String email, required String password}) async {
    final res = await _dio.post(
      '/api/v1/app/login',
      data: {'email': email.trim().toLowerCase(), 'password': password},
    );
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      final token = data is Map ? data['token'] as String? : null;
      if (token != null && token.isNotEmpty) {
        await TokenStore.instance.setToken(token);
        return;
      }
      throw Exception('登录成功但未返回 token');
    }
    throw Exception('登录失败: ${res.statusCode}');
  }
}
