import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:leilao_app/bindings.dart';
import 'package:leilao_app/pages/create_user_page/create_user_view.dart';
import 'package:leilao_app/pages/home_page/leilao_view.dart';
import 'package:leilao_app/pages/login_page/login_view.dart';

class LeilaoApp extends StatelessWidget {
  const LeilaoApp({super.key});

  @override
  Widget build(BuildContext context) {
    return GetMaterialApp(
      initialBinding: LeilaoAppBindings(),
      title: 'Sistema de Leilao Distribuido',
      builder: (context, child) => GestureDetector(
        onTap: () => FocusManager.instance.primaryFocus?.unfocus(),
        child: child,
      ),
      locale: const Locale('pt', 'BR'),
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      initialRoute: '/',
      getPages: [
        GetPage(
          name: '/',
          page: () => LoginView(), binding: LeilaoAppBindings()
        ),
        GetPage(
            name: '/home',
            page: () => LeilaoView(),
            binding: LeilaoAppBindings(),
            ),
        GetPage(
            name: '/create-user',
            page: () => const CreateUserUser(),
            binding: LeilaoAppBindings(),
            ),
      ],
    );
  }
}
