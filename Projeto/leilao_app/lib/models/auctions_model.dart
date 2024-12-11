class Auctions {
  int? id;
  String? produto;
  double? lanceInicial;
  String? dataFinalizacao;
  bool? finalizado;
  String? criador;
  String? vencedor;
  double? lanceFinal;
  double? valorMaximo;
  List<Participante>? participantes;

  Auctions({
    this.id,
    this.produto,
    this.lanceInicial,
    this.dataFinalizacao,
    this.finalizado,
    this.criador,
    this.vencedor,
    this.lanceFinal,
    this.valorMaximo,
    this.participantes,
  });

  Auctions.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    produto = json['produto'];
    lanceInicial = json['lanceInicial']?.toDouble();
    dataFinalizacao = json['dataFinalizacao'];
    finalizado = json['finalizado'];
    criador = json['criador'];
    vencedor = json['vencedor'];
    lanceFinal = json['lanceFinal']?.toDouble();
    valorMaximo = json['valorMaximo']?.toDouble();
    if (json['participantes'] != null) {
      participantes = <Participante>[];
      json['participantes'].forEach((v) {
        participantes!.add(Participante.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['produto'] = this.produto;
    data['lanceInicial'] = this.lanceInicial;
    data['dataFinalizacao'] = this.dataFinalizacao;
    data['finalizado'] = this.finalizado;
    data['criador'] = this.criador;
    data['vencedor'] = this.vencedor;
    data['lanceFinal'] = this.lanceFinal;
    data['valorMaximo'] = this.valorMaximo;
    if (this.participantes != null) {
      data['participantes'] =
          this.participantes!.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class Participante {
  String? usuarioEmail;
  double? lance;

  Participante({this.usuarioEmail, this.lance});

  Participante.fromJson(Map<String, dynamic> json) {
    usuarioEmail = json['usuarioEmail'];
    lance = json['lance']?.toDouble();
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['usuarioEmail'] = this.usuarioEmail;
    data['lance'] = this.lance;
    return data;
  }
}
