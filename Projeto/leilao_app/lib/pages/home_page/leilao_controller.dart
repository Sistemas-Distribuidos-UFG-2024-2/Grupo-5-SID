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

  final TextEditingController lanceController = TextEditingController();

  var leiloes = List<Auctions>.empty().obs;
  var leilaoSelecionado = Auctions().obs;
  var isEmailEditable = false.obs;

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
    await _leilaoRepository.createLeilao(data);
    getLeiloes();
  }

  Future<void> getLeilaoById(int id) async {
    leilaoSelecionado.value = await _leilaoRepository.getLeilaoById(id);
  }

  Future<void> loadEmail() async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';
    emailController.text = email;
    update();
  }

  Future<void> saveEmail(String email) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('userEmail', email);
  }

  void toggleEmailEditable() {
    isEmailEditable.value = !isEmailEditable.value;
  }

  Future<void> placeBid(int leilaoId, double amount) async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';

    try {
      await _leilaoRepository.placeBid(leilaoId, email, amount);
      await getLeilaoById(leilaoId);
      Get.snackbar('Sucesso', 'Seu lance foi registrado com sucesso!');
    } catch (e) {
      Get.snackbar('Erro', 'Falha ao registrar o lance.');
      Get.dialog(
        AlertDialog(
          title: Text('Erro'),
          content:
              Text('Falha ao registrar o lance. Por favor, tente novamente.'),
          actions: <Widget>[
            TextButton(
              child: Text('Fechar'),
              onPressed: () {
                Get.back();
              },
            ),
          ],
        ),
      );
    }
  }

  Future<void> inscreverLeilao(int leilaoId, double lance) async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';

    var data = {
      "usuarioEmail": email,
      "lance": lance,
    };

    try {
      await _leilaoRepository.inscreverLeilao(leilaoId, data);
      Get.snackbar('Sucesso', 'Você está participando do leilão!');
    } catch (e) {
      Get.snackbar('Erro', 'Falha ao inscrever no leilão.');
      Get.dialog(
        AlertDialog(
          title: Text('Erro'),
          content:
              Text('Falha ao inscrever no leilão. Por favor, tente novamente.'),
          actions: <Widget>[
            TextButton(
              child: Text('Fechar'),
              onPressed: () {
                Get.back();
              },
            ),
          ],
        ),
      );
    }
    getLeiloes();
  }
}
