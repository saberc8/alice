import 'package:flutter/material.dart';

// Global theme management with primary color #57be6a
class AppTheme {
  static const Color primary = Color(0xFF57BE6A);
  static const Color primaryDark = Color(0xFF2E9B47);
  static const Color primaryLight = Color(0xFFA6E0B7);

  static ThemeData light() {
    final base = ThemeData.light(useMaterial3: true);
    return base.copyWith(
      colorScheme: ColorScheme.fromSeed(seedColor: primary),
      primaryColor: primary,
      appBarTheme: const AppBarTheme(
        backgroundColor: primary,
        foregroundColor: Colors.white,
        elevation: 0,
      ),
      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: primary,
        foregroundColor: Colors.white,
      ),
      tabBarTheme: const TabBarTheme(
        indicatorColor: Colors.white,
        labelColor: Colors.white,
        unselectedLabelColor: Colors.white70,
      ),
      bottomNavigationBarTheme: const BottomNavigationBarThemeData(
        selectedItemColor: primary,
        unselectedItemColor: Colors.grey,
        type: BottomNavigationBarType.fixed,
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
