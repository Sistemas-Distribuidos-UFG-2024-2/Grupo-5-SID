package com.example;

import org.apache.xmlrpc.client.XmlRpcClient;
import org.apache.xmlrpc.client.XmlRpcClientConfigImpl;

import java.net.URL;
import java.util.HashMap;

public class App {
    public static void main(String[] args) {
        try {
            // Configuração do cliente RPC
            XmlRpcClientConfigImpl config = new XmlRpcClientConfigImpl();
            config.setServerURL(new URL("http://localhost:1234"));  // URL do servidor Go

            XmlRpcClient client = new XmlRpcClient();
            client.setConfig(config);

            // Definir os dados da pessoa
            HashMap<String, Object> pessoa = new HashMap<>();
            pessoa.put("Nome", "Júlio");
            pessoa.put("Sexo", "M");  // M para masculino, F para feminino
            pessoa.put("Idade", 20);

            // Chamar o método remoto
            Object[] params = new Object[]{pessoa};
            String resultado = (String) client.execute("ServicoMaioridade.VerificaMaioridade", params);

            // Exibir o resultado
            System.out.println("Resultado: " + resultado);
        } catch (Exception e) {
            System.err.println("Erro: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
