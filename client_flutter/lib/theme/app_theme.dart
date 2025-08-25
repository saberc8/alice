import 'package:flutter/material.dart';
import 'package:client_flutter/ui/we_colors.dart';

// Global theme management with primary color #57be6a
class AppTheme {
  // Switch to WeChat brand color
  static const Color primary = WeColors.brand;
  static const Color primaryDark = WeColors.brandDark;
  static const Color primaryLight = WeColors.brandLight;

  static ThemeData light() {
    final base = ThemeData.light(useMaterial3: true);
    final scheme = ColorScheme.fromSeed(seedColor: primary).copyWith(
      background: Colors.white,
      surface: Colors.white,
      surfaceVariant: Colors.white,
    );
    return base.copyWith(
      colorScheme: scheme,
      primaryColor: primary,
      scaffoldBackgroundColor: Colors.white,
      canvasColor: Colors.white,
      dialogBackgroundColor: Colors.white,
      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: Colors.white,
        surfaceTintColor: Colors.transparent,
        elevation: 8,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(12)),
        ),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.white,
        foregroundColor: WeColors.textPrimary,
        elevation: 0.4,
        surfaceTintColor: Colors.transparent,
        centerTitle: true,
      ),
      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: primary,
        foregroundColor: Colors.white,
      ),
      tabBarTheme: const TabBarTheme(
        indicatorColor: primary,
        labelColor: primary,
        unselectedLabelColor: WeColors.textSecondary,
      ),
      bottomNavigationBarTheme: const BottomNavigationBarThemeData(
        selectedItemColor: primary,
        unselectedItemColor: Colors.grey,
        type: BottomNavigationBarType.fixed,
        backgroundColor: Colors.white,
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
        focusedBorder: OutlineInputBorder(
          borderSide: const BorderSide(color: primary, width: 1.6),
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    );
  }
}
