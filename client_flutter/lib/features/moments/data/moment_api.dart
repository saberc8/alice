import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/dio_client.dart';
import 'moment_models.dart';

class MomentApi {
  final Dio _dio = DioClient().dio;
  Dio get dio => _dio;

  Future<MomentItem> publish({
    required String content,
    List<String> images = const [],
  }) async {
    final resp = await _dio.post(
      '/api/v1/app/moments',
      data: {'content': content, 'images': images},
    );
    final data = resp.data['data'];
    return MomentItem.fromJson(data as Map<String, dynamic>);
  }

  Future<MomentListResponse> listAll({int page = 1, int pageSize = 20}) async {
    final resp = await _dio.get(
      '/api/v1/app/moments',
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    final data = resp.data['data'];
    return MomentListResponse.fromJson(data as Map<String, dynamic>);
  }

  Future<MomentListResponse> listByUser(
    int userId, {
    int page = 1,
    int pageSize = 20,
  }) async {
    final resp = await _dio.get(
      '/api/v1/app/users/$userId/moments',
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    final data = resp.data['data'];
    return MomentListResponse.fromJson(data as Map<String, dynamic>);
  }

  Future<void> delete(int momentId) async {
    await _dio.delete('/api/v1/app/moments/$momentId');
  }
}
