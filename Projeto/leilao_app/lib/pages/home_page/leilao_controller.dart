import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get/get_state_manager/src/simple/get_controllers.dart';
import 'package:leilao_app/models/LeilaoModel.dart';
import 'package:leilao_app/models/auctions_model.dart';
import 'package:leilao_app/repositories/leilao_repository.dart';
import 'package:shared_preferences/shared_preferences.dart';

class LeilaoController extends GetxController {
  final LeilaoRepository _leilaoRepository = Get.find<LeilaoRepository>();

  final TextEditingController emailController = TextEditingController();

  final TextEditingController nomeController = TextEditingController();
  final TextEditingController valorInicial = TextEditingController();
  final TextEditingController dataController = TextEditingController();
  final TextEditingController valorMaximo = TextEditingController();

  var leiloes = List<Auctions>.empty().obs;
  var leilaoSelecionado = Auctions().obs;

  @override
  void onInit() {
    super.onInit();
  }

  void getLeiloes() async {
    leiloes.value = await _leilaoRepository.getLeiloes();
  }

  void createLeilao() async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';

    var data = {
      "produto": nomeController.text,
      "lanceInicial": int.tryParse(valorInicial.text) ?? 0,
      "dataFinalizacao": dataController.text,
      "criador": email,
      "valorMaximo": int.tryParse(valorMaximo.text) ?? 0
    };
    // await _leilaoRepository.createLeilao(data);
    getLeiloes();
  }

  Future<void> getLeilaoById(int id) async {
    leilaoSelecionado.value = await _leilaoRepository.getLeilaoById(id);
  }

  Future<void> _loadEmail() async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';
    Get.find<LeilaoController>().emailController.text = email;
  }
}
