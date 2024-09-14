package hello_world;

import java.io.*;
import java.net.*;

// Para executar o servidor utilize o comando `java Server.java {PORTA}` exemplo: java Server.java 5601
public class Server {
    public static final int VERIFICATION_PORT = 5610;
    public static final String VERIFICATION_PREFIX = "/health";

    public static final String IP = "localhost";

    private static ServerInfo ServerInfo = new ServerInfo();

    public static void main(String[] args) throws IOException {
        if (args.length == 0) {
            System.out.println("Por favor, forneça a porta como argumento.");
            return;
        }

        int port;
        try {
            port = Integer.parseInt(args[0]);
        } catch (NumberFormatException e) {
            System.out.println("Porta inválida. Use um número inteiro.");
            return;
        }

        ServerInfo = new ServerInfo(IP, port, true);

        ServerSocket serverSocket = new ServerSocket(port);
        System.out.println("Servidor rodando na porta " + port);

        try {
            Socket verificationSocket = new Socket(IP, VERIFICATION_PORT);
            PrintWriter verificationOut = new PrintWriter(verificationSocket.getOutputStream(), true);
            BufferedReader verificationIn = new BufferedReader(new InputStreamReader(verificationSocket.getInputStream()));

            verificationOut.println(VERIFICATION_PREFIX + "/" + ServerInfo.toJSON());
            String verificationResponse = verificationIn.readLine();
            System.out.println("Resposta do servidor de verificação: " + verificationResponse);
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }

        while (true) {
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
}

