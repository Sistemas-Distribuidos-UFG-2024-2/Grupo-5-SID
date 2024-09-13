import java.io.*;
import java.net.*;

public class Server1 {
    public static void main(String[] args) throws IOException {
        ServerSocket serverSocket = new ServerSocket(5000);
        System.out.println("Servidor 1 na porta 5000");


        try (Socket clientSocket = serverSocket.accept()) {
            BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
            PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);

            String message = in.readLine();
            if ("hello".equalsIgnoreCase(message)) {
                out.println("world");
                System.out.println("Mensagem recebida: " + message);
            }
        } catch (IOException e) {
            System.out.println("Erro de conexão no Servidor 1: " + e.getMessage());
        }

    }
}
