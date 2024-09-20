package hello_world;

import java.io.*;
import java.net.*;

public class Client {
    public static final int LOAD_BALANCER_PORT = 5611; // Porta do balanceador
    public static final String IP = "localhost";

    public static void main(String[] args) {
        BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));

        while (true) {
            try (Socket socket = new Socket(IP, LOAD_BALANCER_PORT);
                 PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
                 BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))) {

                // Coleta dados do funcionário
                System.out.print("Nome: ");
                String nome = reader.readLine();
                System.out.print("Cargo: ");
                String cargo = reader.readLine();
                System.out.print("Salário: ");
                double salario = Double.parseDouble(reader.readLine());

                // Envia dados para o servidor
                String message = String.format("FUNCIONARIO,%s,%s,%.2f", nome, cargo, salario);
                out.println(message);

                // Recebe resposta
                String response = in.readLine();
                if (response != null) {
                    System.out.println("Resposta do servidor: " + response);
                } else {
                    System.out.println("Nenhuma resposta recebida.");
                }

                Thread.sleep(3000);
            } catch (IOException | InterruptedException e) {
                System.out.println("Erro na comunicação com o servidor: " + e.getMessage());
            }
        }
    }
}