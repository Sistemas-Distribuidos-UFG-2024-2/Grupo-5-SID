class Auctions {
  int? id;
  String? produto;
  int? lanceInicial;
  String? dataFinalizacao;
  bool? finalizado;
  List<int>? participantes;

  Auctions(
      {this.id,
      this.produto,
      this.lanceInicial,
      this.dataFinalizacao,
      this.finalizado,
      this.participantes});

  Auctions.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    produto = json['produto'];
    lanceInicial = json['lanceInicial'];
    dataFinalizacao = json['dataFinalizacao'];
    finalizado = json['finalizado'];
    participantes = json['participantes'].cast<int>();
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['produto'] = this.produto;
    data['lanceInicial'] = this.lanceInicial;
    data['dataFinalizacao'] = this.dataFinalizacao;
    data['finalizado'] = this.finalizado;
    data['participantes'] = this.participantes;
    return data;
  }
}
