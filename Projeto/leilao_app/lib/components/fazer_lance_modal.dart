import 'package:flutter/material.dart';

class LanceDialog extends StatelessWidget {
  final String nome;
  final double valorAtual;
  final String data;
  final Function(double) onConfirm;

  LanceDialog({
    required this.nome,
    required this.valorAtual,
    required this.data,
    required this.onConfirm,
  });

  @override
  Widget build(BuildContext context) {
    TextEditingController lanceController = TextEditingController();

    return AlertDialog(
      title: Text('Faça seu lance para: $nome'),
      content: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Text('Valor Atual: R\$ ${valorAtual.toStringAsFixed(2)}'),
          Text('Data: $data'),
          SizedBox(height: 16.0),
          TextField(
            controller: lanceController,
            decoration: InputDecoration(
              labelText: 'Valor do Lance',
              border: OutlineInputBorder(),
            ),
            keyboardType: TextInputType.number,
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.of(context).pop();
          },
          child: Text('Cancelar'),
        ),
        ElevatedButton(
          onPressed: () {
            double novoLance = double.parse(lanceController.text);
            onConfirm(novoLance);
            Navigator.of(context).pop();
          },
          child: Text('Confirmar'),
        ),
      ],
    );
  }
}