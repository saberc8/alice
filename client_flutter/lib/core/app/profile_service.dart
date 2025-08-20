import 'dart:io';
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

  Future<Map<String, dynamic>> updateProfile({
    String? nickname,
    String? gender,
    String? bio,
  }) async {
    final body = <String, dynamic>{};
    if (nickname != null) body['nickname'] = nickname;
    if (gender != null) body['gender'] = gender;
    if (bio != null) body['bio'] = bio;
    final res = await _dio.put('/api/v1/app/profile', data: body);
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      if (data is Map<String, dynamic>) return data;
      throw Exception('响应格式错误');
    }
    throw Exception('更新资料失败: ${res.statusCode}');
  }

  Future<Map<String, dynamic>> uploadAvatar(File file) async {
    final form = FormData.fromMap({
      'file': await MultipartFile.fromFile(
        file.path,
        filename: file.path.split('/').last,
      ),
    });
    final res = await _dio.post('/api/v1/app/profile/avatar', data: form);
    if (res.statusCode == 200) {
      final data = res.data is Map ? res.data['data'] : null;
      if (data is Map<String, dynamic>) return data;
      throw Exception('响应格式错误');
    }
    throw Exception('上传头像失败: ${res.statusCode}');
  }

  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    // 后端暂未提供修改密码接口，如后端实现后替换路径与字段
    final res = await _dio.post(
      '/api/v1/app/profile/password',
      data: {'old_password': oldPassword, 'new_password': newPassword},
    );
    if (res.statusCode != 200) {
      throw Exception('修改密码失败: ${res.statusCode}');
    }
  }
}
