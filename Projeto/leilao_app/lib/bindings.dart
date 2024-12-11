import 'package:get/get.dart';
import 'package:leilao_app/pages/home_page/leilao_controller.dart';
import 'package:leilao_app/repositories/leilao_repository.dart';

class LeilaoAppBindings extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<LeilaoRepository>(() => LeilaoRepository());
    Get.lazyPut<LeilaoController>(() => LeilaoController());
  }
}