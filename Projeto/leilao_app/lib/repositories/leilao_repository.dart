import 'package:dio/dio.dart';
import 'package:leilao_app/models/auctions_model.dart';

class LeilaoRepository {
  final dio = Dio();
  

  LeilaoRepository();
  final String baseUrl = 'http://localhost:8080';
  final String baseUrlLance = 'http://auction-app.localhost:6003';

 Future<List<Auctions>> getLeiloes() async {
  try {
    final response = await dio.get('$baseUrl/auctions');
    return (response.data as List)
        .map((leilao) => Auctions.fromJson(leilao))
        .toList();
  } catch (e) {
    print('Erro ao carregar leiloes: $e');
    throw Exception('Failed to load leiloes: $e');
  }
}

   Future<Response> createLeilao(Map<String, dynamic> data) async {
    try {
      final response = await dio.post('$baseUrl/auctions', data: data);
      print('Resposta do servidor: ${response.data}');
      print('Status code: ${response.statusCode}');
      print('Headers: ${response.headers}');
      return response;
    } catch (e) {
      if (e is DioError) {
        print('Erro ao criar leilao: ${e.message}');
        print('Detalhes do erro: ${e.response?.data}');
        print('Status code: ${e.response?.statusCode}');
        print('Headers: ${e.response?.headers}');
      }
      throw Exception('Failed to create leilao: $e');
    }
  }




  Future<Auctions> getLeilaoById(int id) async {
    try {
      final response = await dio.get('$baseUrl/auctions/$id');
      return Auctions.fromJson(response.data);
    } catch (e) {
      throw Exception('Failed to load leilao: $e');
    }
  }

  Future<Response> placeBid(
      int auctionId, String accountEmail, double amount) async {
    try {
      final response = await dio.post(
        '$baseUrlLance/auctions/$auctionId/bids',
        data: {
          'account_email': accountEmail,
          'amount': amount,
        },
      );
      print('Resposta do servidor: ${response.data}');
      print('Status code: ${response.statusCode}');
      print('Headers: ${response.headers}');
      return response;
    } catch (e) {
      if (e is DioError) {
        print('Erro ao dar lance: ${e.message}');
        print('Detalhes do erro: ${e.response?.data}');
        print('Status code: ${e.response?.statusCode}');
        print('Headers: ${e.response?.headers}');
      }
      throw Exception('Failed to place bid: $e');
    }
  }

  Future<Response> inscreverLeilao(int id, Map<String, dynamic> data) async {
    try {
      final response =
          await dio.post('$baseUrl/auctions/$id/inscrever', data: data);
      return response;
    } catch (e) {
      throw Exception('Failed to subscribe to leilao: $e');
    }
  }

  // Future<Response> updateLeilao(int id, Map<String, dynamic> data) async {
  //   try {
  //     final response = await _dio.put('/leiloes/$id', data: data);
  //     return response;
  //   } catch (e) {
  //     throw Exception('Failed to update leilao: $e');
  //   }
  // }

  // Future<Response> deleteLeilao(int id) async {
  //   try {
  //     final response = await _dio.delete('/leiloes/$id');
  //     return response;
  //   } catch (e) {
  //     throw Exception('Failed to delete leilao: $e');
  //   }
  // }
}
