import 'package:flutter/material.dart';

class FazerLanceModal extends StatefulWidget {
  final Function(double) onLanceFeito;

  FazerLanceModal({required this.onLanceFeito});

  @override
  _FazerLanceModalState createState() => _FazerLanceModalState();
}

class _FazerLanceModalState extends State<FazerLanceModal> {
  final _formKey = GlobalKey<FormState>();
  final _lanceController = TextEditingController();

  @override
  void dispose() {
    _lanceController.dispose();
    super.dispose();
  }

  void _submitLance() {
    if (_formKey.currentState!.validate()) {
      final lance = double.parse(_lanceController.text);
      widget.onLanceFeito(lance);
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text('Fazer Lance'),
      content: Form(
        key: _formKey,
        child: TextFormField(
          controller: _lanceController,
          keyboardType: TextInputType.numberWithOptions(decimal: true),
          decoration: InputDecoration(labelText: 'Valor do Lance'),
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
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text('Cancelar'),
        ),
        ElevatedButton(
          onPressed: _submitLance,
          child: Text('Fazer Lance'),
        ),
      ],
    );
  }
}
