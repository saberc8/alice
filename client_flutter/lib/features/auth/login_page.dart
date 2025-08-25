import 'dart:math';
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:client_flutter/theme/app_theme.dart';
import 'package:client_flutter/core/auth/auth_service.dart';
import 'package:flutter/scheduler.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key, required this.onLogin});

  final void Function() onLogin;

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _emailCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  final _nicknameCtrl = TextEditingController();
  bool _loading = false;
  bool _isRegister = false;

  @override
  void dispose() {
    _emailCtrl.dispose();
    _pwdCtrl.dispose();
    _nicknameCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _loading = true);
    try {
      if (_isRegister) {
        await AuthService().register(
          email: _emailCtrl.text,
          password: _pwdCtrl.text,
          nickname: _nicknameCtrl.text,
        );
      } else {
        await AuthService().login(
          email: _emailCtrl.text,
          password: _pwdCtrl.text,
        );
      }
      if (!mounted) return;
      widget.onLogin();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${_isRegister ? '注册' : '登录'}失败: $e')),
      );
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final viewInsets = MediaQuery.of(context).viewInsets;
    return Scaffold(
      resizeToAvoidBottomInset: false,
      body: Stack(
        children: [
          const _AnimatedBlurBackground(),
          // 全屏高斯模糊层：让输入区域的玻璃模糊扩散到整个屏幕
          Positioned.fill(
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 28, sigmaY: 28),
              child: Container(
                // 轻微白蒙版，增强可读性又保留色彩
                color: Colors.white.withOpacity(0.12),
              ),
            ),
          ),
          AnimatedPadding(
            duration: const Duration(milliseconds: 250),
            padding: EdgeInsets.only(bottom: viewInsets.bottom),
            child: Center(
              child: SingleChildScrollView(
                physics: const BouncingScrollPhysics(),
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 380),
                  child: _buildFormCard(context),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFormCard(BuildContext context) {
    // Light style form card
    return Form(
      key: _formKey,
      child: Theme(
        data: Theme.of(context).copyWith(
          inputDecorationTheme: InputDecorationTheme(
            filled: true,
            fillColor: Colors.white.withOpacity(0.72),
            labelStyle: TextStyle(color: Colors.black87.withOpacity(0.72)),
            floatingLabelStyle: const TextStyle(color: Colors.black87),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: BorderSide(color: Colors.black12.withOpacity(0.15)),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: BorderSide(color: Colors.black12.withOpacity(0.12)),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: const BorderSide(color: AppTheme.primary),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: const BorderSide(color: Colors.redAccent),
            ),
            contentPadding: const EdgeInsets.symmetric(
              horizontal: 18,
              vertical: 14,
            ),
          ),
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(28),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 28, vertical: 40),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(28),
              color: Colors.white.withOpacity(0.65),
              border: Border.all(
                color: Colors.white.withOpacity(0.85),
                width: 1,
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.05),
                  blurRadius: 26,
                  offset: const Offset(0, 6),
                ),
              ],
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Text(
                  _isRegister ? '创建账号' : '欢迎回来',
                  textAlign: TextAlign.center,
                  style: const TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                    letterSpacing: 0.5,
                  ),
                ),
                const SizedBox(height: 6),
                Text(
                  _isRegister ? '填写信息完成注册' : '请输入账号信息登录',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87.withOpacity(0.55),
                  ),
                ),
                const SizedBox(height: 32),
                TextFormField(
                  controller: _emailCtrl,
                  decoration: const InputDecoration(labelText: '邮箱'),
                  keyboardType: TextInputType.emailAddress,
                  validator: (v) {
                    if (v == null || v.trim().isEmpty) return '请输入邮箱';
                    final email = v.trim();
                    final ok = RegExp(
                      r'^[^@\s]+@[^@\s]+\.[^@\s]+$',
                    ).hasMatch(email);
                    return ok ? null : '邮箱格式不正确';
                  },
                ),
                const SizedBox(height: 16),
                TextFormField(
                  controller: _pwdCtrl,
                  decoration: const InputDecoration(labelText: '密码'),
                  obscureText: true,
                  validator: (v) {
                    if (v == null || v.isEmpty) return '请输入密码';
                    if (_isRegister && v.length < 6) return '密码至少 6 位';
                    return null;
                  },
                ),
                if (_isRegister) ...[
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _nicknameCtrl,
                    decoration: const InputDecoration(labelText: '昵称（可选）'),
                  ),
                ],
                const SizedBox(height: 32),
                SizedBox(
                  height: 52,
                  child: FilledButton(
                    onPressed: _loading ? null : _submit,
                    style: FilledButton.styleFrom(
                      backgroundColor: AppTheme.primary,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      textStyle: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    child:
                        _loading
                            ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                strokeWidth: 2.2,
                                color: Colors.white,
                              ),
                            )
                            : Text(_isRegister ? '注册并登录' : '登录'),
                  ),
                ),
                const SizedBox(height: 18),
                TextButton(
                  onPressed:
                      _loading
                          ? null
                          : () => setState(() => _isRegister = !_isRegister),
                  child: Text(
                    _isRegister ? '已有账号？去登录' : '没有账号？去注册',
                    style: const TextStyle(color: AppTheme.primary),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

// ================= Animated Multicolor Diffuse Background ================= //

class _AnimatedBlurBackground extends StatefulWidget {
  const _AnimatedBlurBackground();

  @override
  State<_AnimatedBlurBackground> createState() =>
      _AnimatedBlurBackgroundState();
}

class _AnimatedBlurBackgroundState extends State<_AnimatedBlurBackground>
    with SingleTickerProviderStateMixin {
  // 使用连续时间的 ticker，避免周期重启导致的视觉“跳变”
  late final Ticker _ticker;
  double _timeSeconds = 0; // 累积的运行秒数
  final Random _rand = Random(42);
  late final List<_BlobConfig> _blobs;

  @override
  void initState() {
    super.initState();
    _blobs = _generateBlobs();
    _ticker = createTicker((elapsed) {
      // elapsed 会一直增长，保持平滑无缝循环（通过三角函数本身的周期性）
      setState(() {
        _timeSeconds = elapsed.inMilliseconds / 1000.0;
      });
    });
    _ticker.start();
  }

  List<_BlobConfig> _generateBlobs() {
    final colors = [
      const Color(0xFF5B8CFF),
      const Color(0xFFFF6EC7),
      const Color(0xFF6D5BFF),
      const Color(0xFF4ADE80),
      const Color(0xFFFFD166),
      const Color(0xFF23A6D5),
    ];
    return List.generate(6, (i) {
      final base = Offset(_rand.nextDouble(), _rand.nextDouble());
      final amp = 0.15 + _rand.nextDouble() * 0.25;
      return _BlobConfig(
        color: colors[i % colors.length],
        base: base,
        radiusFactor: 0.28 + _rand.nextDouble() * 0.25,
        dxAmp: amp,
        dyAmp: amp * (0.6 + _rand.nextDouble() * 0.8),
        dxFreq: 0.4 + _rand.nextDouble() * 0.8,
        dyFreq: 0.4 + _rand.nextDouble() * 0.8,
        phase: _rand.nextDouble() * pi * 2,
        pulseAmp: 0.06 + _rand.nextDouble() * 0.06,
        pulseFreq: 0.25 + _rand.nextDouble() * 0.35,
      );
    });
  }

  @override
  void dispose() {
    _ticker.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Positioned.fill(
      child: DecoratedBox(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            colors: [Color(0xFFF5F9FF), Color(0xFFFFFFFF)],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: RepaintBoundary(
          // 直接使用 _timeSeconds 传给 painter；CustomPainter 只在 setState 后重绘
          child: CustomPaint(
            painter: _BlobPainter.timeBased(_timeSeconds, _blobs),
            child: const SizedBox.expand(),
          ),
        ),
      ),
    );
  }
}

class _BlobConfig {
  _BlobConfig({
    required this.color,
    required this.base,
    required this.radiusFactor,
    required this.dxAmp,
    required this.dyAmp,
    required this.dxFreq,
    required this.dyFreq,
    required this.phase,
    required this.pulseAmp,
    required this.pulseFreq,
  });
  final Color color;
  final Offset base; // 0..1 relative
  final double radiusFactor; // relative to shortest side
  final double dxAmp;
  final double dyAmp;
  final double dxFreq;
  final double dyFreq;
  final double phase;
  final double pulseAmp; // breathing scale
  final double pulseFreq; // relative speed
}

class _BlobPainter extends CustomPainter {
  _BlobPainter(this.t, this.blobs, {this.timeMode = false});
  _BlobPainter.timeBased(double seconds, List<_BlobConfig> blobs)
    : this(seconds, blobs, timeMode: true);

  final double t; // 如果 timeMode=false 则为 0..1；否则为累计秒数
  final List<_BlobConfig> blobs;
  final bool timeMode; // true => t 是秒数

  @override
  void paint(Canvas canvas, Size size) {
    final shortest = min(size.width, size.height);
    final overlayPaint = Paint()..color = Colors.white.withOpacity(0.35);
    canvas.drawRect(Offset.zero & size, overlayPaint);
    // 与之前 18s 一个周期保持接近速度：之前 angle = (elapsed/18)*2π
    final ttGlobal = timeMode ? (t * (2 * pi / 18.0)) : (t * 2 * pi);
    // First diffuse layer
    canvas.saveLayer(
      Offset.zero & size,
      Paint()..imageFilter = ImageFilter.blur(sigmaX: 60, sigmaY: 60),
    );
    for (final b in blobs) {
      final dx = sin(ttGlobal * b.dxFreq + b.phase) * b.dxAmp;
      final dy = cos(ttGlobal * b.dyFreq + b.phase) * b.dyAmp;
      final pulse = 1 + sin(ttGlobal * b.pulseFreq + b.phase) * b.pulseAmp;
      final center = Offset(
        (b.base.dx + dx).clamp(0.0, 1.0) * size.width,
        (b.base.dy + dy).clamp(0.0, 1.0) * size.height,
      );
      final radius = b.radiusFactor * shortest * pulse;
      final rect = Rect.fromCircle(center: center, radius: radius);
      final grad = RadialGradient(
        colors: [b.color.withOpacity(0.55), b.color.withOpacity(0.0)],
      );
      final paint =
          Paint()
            ..shader = grad.createShader(rect)
            ..blendMode = BlendMode.plus;
      canvas.drawCircle(center, radius, paint);
    }
    canvas.restore();
    // Second highlight layer
    for (final b in blobs) {
      final dx = sin(ttGlobal * b.dxFreq + b.phase + pi / 3) * b.dxAmp * 0.8;
      final dy = cos(ttGlobal * b.dyFreq + b.phase + pi / 3) * b.dyAmp * 0.8;
      final pulse =
          1 + sin(ttGlobal * b.pulseFreq + b.phase + pi / 4) * b.pulseAmp * 0.6;
      final center = Offset(
        (b.base.dx + dx).clamp(0.0, 1.0) * size.width,
        (b.base.dy + dy).clamp(0.0, 1.0) * size.height,
      );
      final radius = b.radiusFactor * shortest * 0.75 * pulse;
      final rect = Rect.fromCircle(center: center, radius: radius);
      final grad = RadialGradient(
        colors: [b.color.withOpacity(0.40), b.color.withOpacity(0.0)],
      );
      final paint =
          Paint()
            ..shader = grad.createShader(rect)
            ..blendMode = BlendMode.plus
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 90);
      canvas.drawCircle(center, radius, paint);
    }
  }

  @override
  bool shouldRepaint(covariant _BlobPainter oldDelegate) =>
      oldDelegate.t != t || oldDelegate.timeMode != timeMode;
}
