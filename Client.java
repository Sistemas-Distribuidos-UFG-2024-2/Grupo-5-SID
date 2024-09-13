import java.io.*;
import java.net.*;

public class Client {
    public static void main(String[] args) {
        String[] servers = {"localhost", "localhost"};
        int[] ports = {5000, 5001};

        for (int i = 0; i < servers.length; i++) {
            try (Socket socket = new Socket(servers[i], ports[i]);
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
