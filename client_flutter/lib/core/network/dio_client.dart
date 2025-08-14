import 'package:dio/dio.dart';
// no flutter foundation needed after unifying base URL

class DioClient {
  static final DioClient _instance = DioClient._internal();
  factory DioClient() => _instance;
  DioClient._internal();

  static String _resolveBaseUrl() {
    const env = String.fromEnvironment('API_BASE_URL');
    if (env.isNotEmpty) return env; // allow --dart-define override

    // Unified default API base URL for all platforms
    return 'http://172.20.121.96:8090';
  }

  late final Dio dio = Dio(
      BaseOptions(
        baseUrl: _resolveBaseUrl(),
        connectTimeout: const Duration(seconds: 5),
        receiveTimeout: const Duration(seconds: 10),
        headers: const {'Content-Type': 'application/json'},
      ),
    )
    ..interceptors.addAll([
      LogInterceptor(requestBody: false, responseBody: true),
    ]);
}
