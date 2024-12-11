import 'dart:async';

import 'package:flutter/material.dart';
import 'package:leilao_app/leilao_app.dart';

FutureOr<void> main() async {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return LeilaoApp();
  }
}
