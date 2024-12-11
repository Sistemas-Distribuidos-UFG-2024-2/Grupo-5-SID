import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:leilao_app/pages/home_page/leilao_controller.dart';

class InscreverLeilaoModal extends StatefulWidget {
  final Function(double) onLanceFeito;

  InscreverLeilaoModal({required this.onLanceFeito});

  @override
  _InscreverLeilaoModalState createState() => _InscreverLeilaoModalState();
}

class _InscreverLeilaoModalState extends State<InscreverLeilaoModal> {
  final _formKey = GlobalKey<FormState>();
  final _lanceController = TextEditingController();
  final controller = Get.find<LeilaoController>();

  @override
  void dispose() {
    _lanceController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text('Inscrever no Leilão'),
      content: Container(
        height: 120,
        child: Column(
          children: [
            Text('Para se inscrever no leilão, faça um lance'),
            SizedBox(height: 16),
            Form(
              key: _formKey,
              child: TextFormField(
                controller: _lanceController,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  labelText: 'Valor do Lance',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Por favor, insira um valor';
                  }
                  if (double.tryParse(value) == null) {
                    return 'Por favor, insira um valor válido';
                  }
                  return null;
                },
              ),
            ),
          ],
        ),
      ),
      actions: <Widget>[
        TextButton(
          child: Text('Cancelar'),
          onPressed: () {
            Navigator.of(context).pop();
          },
        ),
        ElevatedButton(
          child: Text('Confirmar'),
          onPressed: () {
            if (_formKey.currentState!.validate()) {
              double lance = double.parse(_lanceController.text);
              widget.onLanceFeito(lance);
              Navigator.of(context).pop();
            }
          },
        ),
      ],
    );
  }
}
