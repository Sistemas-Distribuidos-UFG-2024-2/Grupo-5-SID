package hello_world;

import java.io.*;
import java.net.*;

public class Client {
    public static void main(String[] args) {
        int[] ports = {5601, 5602, 5603};

        while (true) {
            try {
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                System.out.println("Erro no sleep");
            }
            for (int i = 0; i < ports.length; i++) {
                try (Socket socket = new Socket("localhost", ports[i]);
                     PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
                     BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))) {

                    out.println("hello");
                    String response = in.readLine();
                    System.out.println("Resposta do servidor " + socket.getPort() + " : " + response);
                    break;

                } catch (IOException e) {
                    System.out.println("Servidor " + (i + 1) + " indisponível. Tentando próximo...");
                }
            }
        }
    }
}
