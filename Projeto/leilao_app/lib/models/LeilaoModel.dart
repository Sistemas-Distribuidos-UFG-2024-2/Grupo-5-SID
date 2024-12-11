
class Leilao {
  String nome;
  String tipo;
  String data;
  double ultimoLance;

  Leilao(
      {required this.nome,
      required this.tipo,
      required this.data,
      this.ultimoLance = 0.0});
}