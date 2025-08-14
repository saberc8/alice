import 'package:dio/dio.dart';
import 'package:client_flutter/core/network/dio_client.dart';

class HealthService {
  final Dio _dio = DioClient().dio;

  Future<bool> ping() async {
    try {
      final res = await _dio.get('/health');
      if (res.statusCode == 200) {
        final data = res.data;
        if (data is Map && data['status'] == 'ok') return true;
        return true; // tolerant to non-schema ok 200
      }
      return false;
    } on DioException catch (_) {
      return false;
    } catch (_) {
      return false;
    }
  }
}
