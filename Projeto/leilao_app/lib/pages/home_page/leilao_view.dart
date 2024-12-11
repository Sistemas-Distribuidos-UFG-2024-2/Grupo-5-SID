import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';
import 'package:leilao_app/components/fazer_lance_modal.dart';
import 'package:leilao_app/models/LeilaoModel.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'leilao_controller.dart';

class LeilaoView extends StatefulWidget {
  const LeilaoView({super.key});

  @override
  State<LeilaoView> createState() => _LeilaoViewState();
}

class _LeilaoViewState extends State<LeilaoView> {
  @override
  void initState() {
    super.initState();
    _loadEmail();
    Get.find<LeilaoController>().getLeiloes();
  }

  Future<void> _loadEmail() async {
    final prefs = await SharedPreferences.getInstance();
    final email = prefs.getString('userEmail') ?? '';
    Get.find<LeilaoController>().emailController.text = email;
  }

  @override
  Widget build(BuildContext context) {
    return GetBuilder<LeilaoController>(
      init: LeilaoController(),
      builder: (controller) {
        return Scaffold(
          appBar: AppBar(
            title: Text('Sistema de Leilao Distribuido',
                style: TextStyle(fontSize: 32)),
            centerTitle: true,
          ),
          body: Container(
            child: Padding(
              padding: const EdgeInsets.all(32.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: MediaQuery.of(context).size.width * 0.45,
                    child: Column(
                      children: [
                        Text('Criar Leilao', style: TextStyle(fontSize: 24)),
                        Padding(
                          padding: const EdgeInsets.all(16.0),
                          child: Center(
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: <Widget>[
                                TextField(
                                  controller: controller.nomeController,
                                  decoration: InputDecoration(
                                    labelText: 'Nome do Produto',
                                    border: OutlineInputBorder(),
                                  ),
                                ),
                                SizedBox(height: 16),
                                TextField(
                                  controller: controller.valorInicial,
                                  decoration: InputDecoration(
                                    labelText: 'Valor inicial',
                                    border: OutlineInputBorder(),
                                  ),
                                ),
                                SizedBox(height: 16),
                                TextField(
                                  decoration: InputDecoration(
                                    enabled: false,
                                    labelText:
                                        controller.emailController.text.isEmpty
                                            ? 'Email do criador'
                                            : controller.emailController.text,
                                    border: OutlineInputBorder(),
                                  ),
                                ),
                                SizedBox(height: 16),
                                TextField(
                                  controller: controller.dataController,
                                  decoration: InputDecoration(
                                    labelText: 'Data final do leilão',
                                    border: OutlineInputBorder(),
                                  ),
                                  onTap: () async {
                                    DateTime? pickedDate = await showDatePicker(
                                      context: context,
                                      initialDate: DateTime.now(),
                                      firstDate: DateTime(2000),
                                      lastDate: DateTime(2101),
                                    );

                                    if (pickedDate != null) {
                                      TimeOfDay? pickedTime =
                                          await showTimePicker(
                                        context: context,
                                        initialTime: TimeOfDay.fromDateTime(
                                            DateTime.now()),
                                        builder: (BuildContext context,
                                            Widget? child) {
                                          return MediaQuery(
                                            data: MediaQuery.of(context)
                                                .copyWith(
                                                    alwaysUse24HourFormat:
                                                        true),
                                            child: child!,
                                          );
                                        },
                                      );
                                      if (pickedTime != null) {
                                        DateTime finalDateTime = DateTime(
                                          pickedDate.year,
                                          pickedDate.month,
                                          pickedDate.day,
                                          pickedTime.hour,
                                          pickedTime.minute,
                                        );

                                        String formattedDateTime =
                                            DateFormat("yyyy-MM-dd'T'HH:mm:ss")
                                                .format(finalDateTime);

                                        controller.dataController.text =
                                            formattedDateTime;
                                      }
                                    }
                                  },
                                ),
                                SizedBox(height: 32),
                                ElevatedButton(
                                  onPressed: () {
                                    controller.createLeilao();
                                  },
                                  child: Text('Criar'),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  // SizedBox(height: 32),
                  Container(
                    width: MediaQuery.of(context).size.width * 0.45,
                    child: Column(
                      children: [
                        Text('Leilões Ativos', style: TextStyle(fontSize: 24)),
                        Expanded(
                          child: Obx(() {
                            return ListView.builder(
                              itemCount: controller.leiloes.length,
                              itemBuilder: (context, index) {
                                final leilao = controller.leiloes[index];
                                return Card(
                                  child: ListTile(
                                    title: Text(leilao.produto.toString()),
                                    subtitle: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.start,
                                      children: [
                                        Text(
                                            'Lance inicial: R\$ ${leilao.lanceInicial}'),
                                        Text(
                                            'Data Finalização: ${leilao.dataFinalizacao}'),
                                        ('${leilao.finalizado}' == 'true')
                                            ? Text('Status: Finalizado',
                                                style: TextStyle(
                                                    color: Colors.red))
                                            : Text('Status: Ativo',
                                                style: TextStyle(
                                                    color: Colors.green)),
                                      ],
                                    ),
                                    trailing: Column(
                                        mainAxisSize: MainAxisSize.min,
                                        children: [
                                          Expanded(
                                            child: ElevatedButton(
                                              style: ButtonStyle(
                                                  textStyle:
                                                      MaterialStateProperty.all(
                                                    const TextStyle(
                                                      color: Colors.white,
                                                    ),
                                                  ),
                                                  backgroundColor: leilao
                                                              .finalizado ==
                                                          true
                                                      ? MaterialStateProperty
                                                          .all(Colors.grey)
                                                      : MaterialStateProperty
                                                          .all(Colors.white)),
                                              onPressed: () {
                                                if (leilao.finalizado == true) {
                                                  showAboutDialog(
                                                      context: context,
                                                      applicationName:
                                                          'Leilão Finalizado',
                                                      applicationVersion:
                                                          'O leilão já foi finalizado, não é possível fazer lances.');
                                                } else {}
                                              },
                                              child: Text('Fazer Lance'),
                                            ),
                                          ),
                                          SizedBox(height: 4),
                                          Expanded(
                                            child: ElevatedButton(
                                              style: ButtonStyle(
                                                foregroundColor:
                                                    MaterialStateProperty.all(
                                                        Colors.white),
                                                textStyle:
                                                    MaterialStateProperty.all(
                                                  const TextStyle(
                                                      color: Colors.white),
                                                ),
                                                backgroundColor:
                                                    MaterialStateProperty.all(
                                                        Colors.blue),
                                              ),
                                              onPressed: () async {
                                                await controller
                                                    .getLeilaoById(leilao.id!);
                                                showDialog(
                                                  context: context,
                                                  builder:
                                                      (BuildContext context) {
                                                    return AlertDialog(
                                                      title: Text(
                                                          'Detalhes do Leilão'),
                                                      content:
                                                          SingleChildScrollView(
                                                        child: ListBody(
                                                          children: <Widget>[
                                                            Text(
                                                                'Produto: ${controller.leilaoSelecionado?.value.produto ?? ''}'),
                                                            Text(
                                                                'Lance Inicial: R\$ ${controller.leilaoSelecionado?.value.lanceInicial ?? ''}'),
                                                            Text(
                                                                'Data de Finalização: ${controller.leilaoSelecionado?.value.dataFinalizacao ?? ''}'),
                                                            Text(
                                                                'Status: ${controller.leilaoSelecionado?.value.finalizado == true ? 'Finalizado' : 'Ativo'}'),
                                                          ],
                                                        ),
                                                      ),
                                                      actions: <Widget>[
                                                        TextButton(
                                                          child: Text('Fechar'),
                                                          onPressed: () {
                                                            Navigator.of(
                                                                    context)
                                                                .pop();
                                                          },
                                                        ),
                                                      ],
                                                    );
                                                  },
                                                );
                                              },
                                              child: Text('Saiba Mais'),
                                            ),
                                          ),
                                        ]),
                                  ),
                                );
                              },
                            );
                          }),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
