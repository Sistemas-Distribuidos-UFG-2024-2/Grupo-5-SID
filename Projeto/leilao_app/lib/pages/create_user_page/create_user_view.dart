import 'package:flutter/material.dart';
import 'package:get/get_core/src/get_main.dart';

class CreateUserUser extends StatefulWidget {
  const CreateUserUser({super.key});

  @override
  _CreateUserUserState createState() => _CreateUserUserState();
}

class _CreateUserUserState extends State<CreateUserUser> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();

  void _createUser() {
    if (_formKey.currentState!.validate()) {
      // Handle user creation logic here
      String email = emailController.text;
      String password = passwordController.text;
      print('Email: $email, Password: $password');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Criar Conta', style: TextStyle(fontSize: 32)),
      ),
      body: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(8.0),
              color: Colors.grey[200],
            ),
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
                  const TextField(
                    decoration: InputDecoration(
                      labelText: 'Email',
                      border: OutlineInputBorder(),
                    ),
                  ),
                  SizedBox(height: 16.0),
                  const TextField(
                    decoration: InputDecoration(
                      labelText: 'Usuário',
                      border: OutlineInputBorder(),
                    ),
                  ),
                  SizedBox(height: 16.0),
                  ElevatedButton(
                    onPressed: () {
                      // Handle login logic here
                    },
                    child: Text('Criar Conta', style: TextStyle(fontSize: 24)),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
