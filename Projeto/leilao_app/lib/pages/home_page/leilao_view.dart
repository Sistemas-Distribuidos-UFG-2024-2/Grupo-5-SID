import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';
import 'package:leilao_app/components/fazer_lance_modal.dart';
import 'package:leilao_app/models/LeilaoModel.dart';
import 'package:leilao_app/models/inscrever_leilao_modal.dart';
import 'package:leilao_app/models/lance_modal.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'leilao_controller.dart';

class LeilaoView extends StatefulWidget {
  const LeilaoView({super.key});

  @override
  State<LeilaoView> createState() => _LeilaoViewState();
}

class _LeilaoViewState extends State<LeilaoView> {
  // @override
  // void initState() {
  //   super.initState();
  //   Get.find<LeilaoController>().getLeiloes();
  // }

  @override
  Widget build(BuildContext context) {
    return GetBuilder<LeilaoController>(
      init: LeilaoController(),
      initState: (state) {
        Get.find<LeilaoController>().getLeiloes();
        Get.find<LeilaoController>().loadEmail();
      },
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
                    child: SingleChildScrollView(
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
                                    controller: controller.valorMaximo,
                                    decoration: InputDecoration(
                                      labelText: 'Valor maximo',
                                      border: OutlineInputBorder(),
                                    ),
                                  ),
                                  SizedBox(height: 16),
                                  Obx(() => TextField(
                                        controller: controller.emailController,
                                        enabled:
                                            controller.isEmailEditable.value,
                                        decoration: InputDecoration(
                                          labelText: controller
                                                  .emailController.text.isEmpty
                                              ? 'Email do criador'
                                              : controller.emailController.text,
                                          border: OutlineInputBorder(),
                                        ),
                                      )),
                                  SizedBox(height: 16),
                                  ElevatedButton(
                                    onPressed: () {
                                      if (controller.isEmailEditable.value) {
                                        controller.saveEmail(
                                            controller.emailController.text);
                                        controller.toggleEmailEditable();
                                      } else {
                                        controller.toggleEmailEditable();
                                      }
                                    },
                                    child: Obx(() => Text(
                                        controller.isEmailEditable.value
                                            ? 'Salvar Email'
                                            : 'Editar Email')),
                                  ),
                                  SizedBox(height: 16),
                                  TextField(
                                    controller: controller.dataController,
                                    decoration: InputDecoration(
                                      labelText: 'Data final do leilão',
                                      border: OutlineInputBorder(),
                                    ),
                                    onTap: () async {
                                      DateTime? pickedDate =
                                          await showDatePicker(
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

                                          String formattedDateTime = DateFormat(
                                                  "yyyy-MM-dd'T'HH:mm:ss")
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
                              shrinkWrap: true,
                              itemCount: controller.leiloes.length,
                              itemBuilder: (context, index) {
                                final leilao = controller.leiloes[index];
                                return SizedBox(
                                  height: 120,
                                  child: Card(
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
                                                        MaterialStateProperty
                                                            .all(
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
                                                  if (leilao.finalizado ==
                                                      true) {
                                                    showAboutDialog(
                                                        context: context,
                                                        applicationName:
                                                            'Leilão Finalizado',
                                                        applicationVersion:
                                                            'O leilão já foi finalizado, não é possível fazer lances.');
                                                  } else {
                                                    showDialog(
                                                      context: context,
                                                      builder: (BuildContext
                                                          context) {
                                                        return FazerLanceModal(
                                                          onLanceFeito:
                                                              (lance) async {
                                                            final prefs =
                                                                await SharedPreferences
                                                                    .getInstance();
                                                            final email =
                                                                prefs.getString(
                                                                        'userEmail') ??
                                                                    '';
                                                            await controller
                                                                .placeBid(
                                                                    leilao.id!,
                                                                    lance);
                                                            controller
                                                                .getLeiloes();
                                                          },
                                                        );
                                                      },
                                                    );
                                                  }
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
                                                      .getLeilaoById(
                                                          leilao.id!);
                                                  showDialog(
                                                    context: context,
                                                    builder:
                                                        (BuildContext context) {
                                                      // Pré-construir a lista de participantes
                                                      final participantesWidgets = controller
                                                                  .leilaoSelecionado
                                                                  ?.value
                                                                  .participantes !=
                                                              null
                                                          ? controller
                                                              .leilaoSelecionado!
                                                              .value
                                                              .participantes!
                                                              .map(
                                                                  (participante) {
                                                              return Padding(
                                                                padding:
                                                                    const EdgeInsets
                                                                        .only(
                                                                        bottom:
                                                                            12),
                                                                child: ListTile(
                                                                  title: Text(
                                                                      style: TextStyle(
                                                                          fontSize:
                                                                              12,
                                                                          fontWeight: FontWeight
                                                                              .bold),
                                                                      participante
                                                                              .usuarioEmail ??
                                                                          'N/A'),
                                                                  subtitle: Text(
                                                                      style: TextStyle(
                                                                          fontSize:
                                                                              12,
                                                                          fontWeight: FontWeight
                                                                              .bold,
                                                                          color:
                                                                              Colors.green),
                                                                      'Lance: R\$ ${participante.lance ?? 'N/A'}'),
                                                                ),
                                                              );
                                                            }).toList()
                                                          : [
                                                              Text(
                                                                  'Nenhum participante')
                                                            ];

                                                      return AlertDialog(
                                                        title: Text(
                                                            'Detalhes do Leilão'),
                                                        content:
                                                            SingleChildScrollView(
                                                          child: Column(
                                                            crossAxisAlignment:
                                                                CrossAxisAlignment
                                                                    .start,
                                                            mainAxisSize:
                                                                MainAxisSize
                                                                    .min,
                                                            children: <Widget>[
                                                              Text(
                                                                  'Produto: ${controller.leilaoSelecionado?.value.produto ?? ''}'),
                                                              Text(
                                                                  'Lance Inicial: R\$ ${controller.leilaoSelecionado?.value.lanceInicial ?? ''}'),
                                                              Text(
                                                                  'Data de Finalização: ${controller.leilaoSelecionado?.value.dataFinalizacao.toString().substring(0, 10) ?? ''}'),
                                                              Text(
                                                                  'Status: ${controller.leilaoSelecionado?.value.finalizado == true ? 'Finalizado' : 'Ativo'}'),
                                                              Text(
                                                                  'Vencedor: ${controller.leilaoSelecionado?.value.vencedor ?? 'N/A'}'),
                                                              Text(
                                                                  'Lance Final: R\$ ${controller.leilaoSelecionado?.value.lanceFinal ?? 'N/A'}'),
                                                              Text(
                                                                  'Valor Máximo: R\$ ${controller.leilaoSelecionado?.value.valorMaximo ?? 'N/A'}'),
                                                              SizedBox(
                                                                  height: 8),
                                                              Text(
                                                                  'Participantes:'),
                                                              ...participantesWidgets,
                                                            ],
                                                          ),
                                                        ),
                                                        actions: <Widget>[
                                                          TextButton(
                                                            child:
                                                                Text('Fechar'),
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
                                            SizedBox(width: 16),
                                            Expanded(
                                              child: ElevatedButton(
                                                style: ButtonStyle(
                                                  textStyle:
                                                      MaterialStateProperty.all(
                                                    const TextStyle(
                                                        color: Colors.white),
                                                  ),
                                                ),
                                                onPressed: () async {
                                                  showDialog(
                                                      context: context,
                                                      builder: (context) =>
                                                          InscreverLeilaoModal(
                                                            onLanceFeito:
                                                                (lance) async {
                                                              await controller
                                                                  .inscreverLeilao(
                                                                      leilao
                                                                          .id!,
                                                                      lance);
                                                            },
                                                          ));
                                                },
                                                child:
                                                    Text('Inscrever no Leilão'),
                                              ),
                                            ),
                                          ]),
                                    ),
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
