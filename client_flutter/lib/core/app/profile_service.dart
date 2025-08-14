import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/dio_client.dart';

class ProfileService {
  final Dio _dio = DioClient().dio;

  Future<Map<String, dynamic>> getProfile() async {
    final res = await _dio.get('/api/v1/app/profile');
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      if (data is Map<String, dynamic>) return data;
      throw Exception('响应格式错误');
    }
    throw Exception('获取资料失败: ${res.statusCode}');
  }
}
