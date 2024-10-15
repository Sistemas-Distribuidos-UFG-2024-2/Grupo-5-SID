package com.example;

import org.apache.xmlrpc.client.XmlRpcClient;
import org.apache.xmlrpc.client.XmlRpcClientConfigImpl;

import java.net.URL;
import java.util.HashMap;

public class App {
    public static void main(String[] args) {
        try {
            XmlRpcClientConfigImpl config = new XmlRpcClientConfigImpl();
            config.setServerURL(new URL("http://localhost:1234"));

            XmlRpcClient client = new XmlRpcClient();
            client.setConfig(config);


            HashMap<String, Object> pessoa = new HashMap<>();
            pessoa.put("Nome", "Júlio");
            pessoa.put("Sexo", "M");
            pessoa.put("Idade", 20);


            Object[] params = new Object[]{pessoa};
            String resultado = (String) client.execute("ServicoMaioridade.VerificaMaioridade", params);


            System.out.println("Resultado: " + resultado);
        } catch (Exception e) {
            System.err.println("Erro: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
