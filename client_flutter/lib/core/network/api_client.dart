import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'dio_client.dart';

/// 统一的业务异常，包含可选状态码及后端返回的原始数据
class ApiException implements Exception {
  final String message;
  final int? statusCode;
  final dynamic raw;
  ApiException(this.message, {this.statusCode, this.raw});
  @override
  String toString() =>
      'ApiException(statusCode: $statusCode, message: $message)';
}

/// 对后端固定响应结构进行轻封装：{"data": ...}
class ApiClient {
  ApiClient._();
  static final ApiClient instance = ApiClient._();

  final Dio _dio = DioClient().dio;
  Dio get raw => _dio; // 仅在需要原始能力(如 上传 FormData 的自定义处理)时使用

  /// 通用 GET
  Future<T> get<T>(
    String path, {
    Map<String, dynamic>? query,
    T Function(dynamic data)? parser,
  }) async {
    return _request<T>(
      () => _dio.get(path, queryParameters: query),
      parser: parser,
    );
  }

  /// 通用 POST
  Future<T> post<T>(
    String path, {
    dynamic body,
    Map<String, dynamic>? query,
    T Function(dynamic data)? parser,
    Options? options,
  }) async {
    return _request<T>(
      () =>
          _dio.post(path, data: body, queryParameters: query, options: options),
      parser: parser,
    );
  }

  /// 通用 PUT
  Future<T> put<T>(
    String path, {
    dynamic body,
    Map<String, dynamic>? query,
    T Function(dynamic data)? parser,
  }) async {
    return _request<T>(
      () => _dio.put(path, data: body, queryParameters: query),
      parser: parser,
    );
  }

  /// 通用 DELETE
  Future<T> delete<T>(
    String path, {
    dynamic body,
    Map<String, dynamic>? query,
    T Function(dynamic data)? parser,
  }) async {
    return _request<T>(
      () => _dio.delete(path, data: body, queryParameters: query),
      parser: parser,
    );
  }

  Future<T> _request<T>(
    Future<Response<dynamic>> Function() fn, {
    T Function(dynamic data)? parser,
  }) async {
    try {
      final res = await fn();
      final code = res.statusCode ?? 0;
      if (code < 200 || code >= 300) {
        throw ApiException('HTTP $code', statusCode: code, raw: res.data);
      }
      final root = res.data;
      final data = (root is Map) ? root['data'] : null;
      final dynamic payload = data ?? root; // 若后端未包 data 兼容
      if (parser != null) return parser(payload);
      return payload as T;
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      final raw = e.response?.data;
      // 读取后端可能的 message 字段
      String? backendMsg;
      final respData = raw;
      if (respData is Map) {
        backendMsg =
            respData['message']?.toString() ?? respData['error']?.toString();
      }
      throw ApiException(
        backendMsg ?? e.message ?? '网络请求失败',
        statusCode: status,
        raw: raw,
      );
    } catch (e, s) {
      if (kDebugMode) {
        // ignore: avoid_print
        print('ApiClient 未捕获异常: $e\n$s');
      }
      if (e is ApiException) rethrow;
      throw ApiException(e.toString());
    }
  }
}
