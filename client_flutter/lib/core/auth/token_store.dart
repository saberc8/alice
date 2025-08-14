import 'package:shared_preferences/shared_preferences.dart';

class TokenStore {
  TokenStore._();
  static final TokenStore instance = TokenStore._();

  static const _kKey = 'auth_token';
  String? _token;

  String? get token => _token;

  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    _token = prefs.getString(_kKey);
  }

  Future<void> setToken(String token) async {
    _token = token;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_kKey, token);
  }

  Future<void> clear() async {
    _token = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_kKey);
  }
}
