import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:get/get.dart';
import 'package:leilao_app/pages/home_page/leilao_controller.dart';
import 'package:shared_preferences/shared_preferences.dart';
class LoginView extends StatefulWidget {
  const LoginView({super.key});
  

  @override
  State<LoginView> createState() => _LoginViewState();
}

class _LoginViewState extends State<LoginView> {
  @override
  void initState() {
    super.initState();
    
    Get.find<LeilaoController>().getLeiloes();
  }

  final TextEditingController emailController = TextEditingController();


  @override
  Widget build(BuildContext context) {
    return GetBuilder<LeilaoController>(
      init: LeilaoController(),
      builder: (controller) {
        return Scaffold(
      appBar: AppBar(
        title: Text('Sistema de Leilao Distribuido', style: TextStyle(fontSize: 32)),
      ),
      body: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(8.0),
              color: Colors.grey[200],),
            width: MediaQuery.of(context).size.width * 0.45,
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  Image.asset(
                    'assets/images/leilao.png',
                    height: 100,
                  ),
                  SizedBox(height: 32.0),
                  TextField(
                    controller: emailController,
                    decoration: InputDecoration(
                      labelText: 'email',
                      border: OutlineInputBorder(),
                    ),
                  ),
                  SizedBox(height: 16.0),
            
                  SizedBox(height: 16.0),
                  ElevatedButton(
                    onPressed: () async {
                      // Handle login logic here
                          await controller.saveEmail(emailController.text);
                      Get.toNamed('/home');
                    },
                    child: Text('Login', style: TextStyle(fontSize: 24)),
                  ),

                  SizedBox(height: 16.0),
                      // ElevatedButton(
                      //   onPressed: () {
                      //     // Handle login logic here
                      //     Get.toNamed('/create-user');
                      //   },
                      //   child: Text('Registrar', style: TextStyle(fontSize: 24)),
                      // ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
      },
    );
  }
}
