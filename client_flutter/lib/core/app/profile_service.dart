// ignore: depend_on_referenced_packages
import 'dart:io' show File; // 条件使用，仅在非 Web 环境有效
import 'dart:typed_data';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/api_client.dart';

class ProfileService {
  final _api = ApiClient.instance;
  Dio get _dio => _api.raw; // 暴露底层用于文件上传（仍复用拦截器）

  Future<Map<String, dynamic>> getProfile() async {
    return _api.get<Map<String, dynamic>>(
      '/api/v1/app/profile',
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
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
    return _api.put<Map<String, dynamic>>(
      '/api/v1/app/profile',
      body: body,
      parser: (d) => (d is Map<String, dynamic>) ? d : <String, dynamic>{},
    );
  }

  /// Web 与 移动端统一的上传接口
  /// [file] 仅在非 Web 使用；Web 需传入 [bytes]
  Future<Map<String, dynamic>> uploadAvatar({
    File? file,
    Uint8List? bytes,
    required String filename,
  }) async {
    MultipartFile mf;
    if (kIsWeb) {
      if (bytes == null) throw Exception('Web 上传缺少字节数据');
      mf = MultipartFile.fromBytes(bytes, filename: filename);
    } else {
      final f = file;
      if (f == null) throw Exception('缺少文件');
      mf = await MultipartFile.fromFile(f.path, filename: filename);
    }
    final form = FormData.fromMap({'file': mf});
    final res = await _dio.post('/api/v1/app/profile/avatar', data: form);
    final data = res.data is Map ? res.data['data'] : null;
    if (data is Map<String, dynamic>) return data;
    throw ApiException('响应格式错误');
  }

  Future<void> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    // 后端暂未提供修改密码接口，如后端实现后替换路径与字段
    await _api.post(
      '/api/v1/app/profile/password',
      body: {'old_password': oldPassword, 'new_password': newPassword},
    );
  }
}
