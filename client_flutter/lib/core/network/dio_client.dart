import 'package:dio/dio.dart';
import 'package:client_flutter/core/auth/token_store.dart';
import 'package:flutter/foundation.dart';

class DioClient {
  static final DioClient _instance = DioClient._internal();
  factory DioClient() => _instance;
  DioClient._internal();

  static String _resolveBaseUrl() {
    // 1. Highest priority: compile-time override
    const env = String.fromEnvironment('API_BASE_URL');
    if (env.isNotEmpty)
      return env; // pass with: flutter run --dart-define API_BASE_URL=http://192.168.x.x:8090
    // 2. Platform defaults
    // Android 模拟器访问宿主机 localhost 需要用 10.0.2.2
    // (如果是真机请使用局域网 IP，并通过 --dart-define 传入)
    if (!kIsWeb && defaultTargetPlatform == TargetPlatform.android) {
      return 'http://10.0.2.2:8090';
    }

    // 3. 其它平台本地运行直接 localhost
    return 'http://localhost:8090';
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
      InterceptorsWrapper(
        onRequest: (options, handler) {
          final token = TokenStore.instance.token;
          if (token != null && token.isNotEmpty) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          handler.next(options);
        },
      ),
      LogInterceptor(requestBody: false, responseBody: true),
    ]);
}
