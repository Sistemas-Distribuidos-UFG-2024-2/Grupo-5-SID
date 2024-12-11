// import 'package:flutter/material.dart';
// import 'package:get/get.dart';
// import 'package:shared_preferences/shared_preferences.dart';

// class LoginView extends StatelessWidget {


//   @override
//   Widget build(BuildContext context) {
//     return Scaffold(
//       appBar: AppBar(
//         title: Text('Sistema de Leilao Distribuido', style: TextStyle(fontSize: 32)),
//       ),
//       body: Row(
//         mainAxisAlignment: MainAxisAlignment.center,
//         children: [
//           Container(
//             decoration: BoxDecoration(
//               borderRadius: BorderRadius.circular(8.0),
//               color: Colors.grey[200],),
//             width: MediaQuery.of(context).size.width * 0.45,
//             child: Padding(
//               padding: const EdgeInsets.all(16.0),
//               child: Column(
//                 mainAxisAlignment: MainAxisAlignment.center,
//                 children: <Widget>[
//                   Image.asset(
//                     'assets/images/leilao.png',
//                     height: 100,
//                   ),
//                   SizedBox(height: 32.0),
//                   const TextField(
//                     decoration: InputDecoration(
//                       labelText: 'Username',
//                       border: OutlineInputBorder(),
//                     ),
//                   ),
//                   SizedBox(height: 16.0),
            
//                   SizedBox(height: 16.0),
//                   ElevatedButton(
//                     onPressed: () {
//                       // Handle login logic here
//                       _saveEmail(controller.emailController.text);

//                       Get.toNamed('/home');
//                     },
//                     child: Text('Login', style: TextStyle(fontSize: 24)),
//                   ),

//                   SizedBox(height: 16.0),
//                   ElevatedButton(
//                     onPressed: () {
//                       // Handle login logic here
//                       Get.toNamed('/create-user');
//                     },
//                     child: Text('Registrar', style: TextStyle(fontSize: 24)),
//                   ),
//                 ],
//               ),
//             ),
//           ),
//         ],
//       ),
//     );
//   }
// }