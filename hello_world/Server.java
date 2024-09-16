package hello_world;

import java.io.*;
import java.net.*;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

// Para executar o servidor utilize o comando `java Server.java {PORTA}` exemplo: java Server.java 5601
public class Server {
    public static final int VERIFICATION_PORT = 5610;
    public static final String VERIFICATION_PREFIX = "/health";

    public static final String IP = "localhost";

    private static ServerInfo serverInfo = new ServerInfo();
    private static final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
    public static final int HealthCheckTimeSeconds = 10;

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

        serverInfo = new ServerInfo(IP, port, true);

        ServerSocket serverSocket = new ServerSocket(port);
        System.out.println("Servidor rodando na porta " + port);

        scheduler.scheduleAtFixedRate(() -> {
            try {
                sendServerInfo();
            } catch (IOException e) {
                System.out.println("Erro ao enviar informações do servidor: " + e.getMessage());
            }
        }, 0, HealthCheckTimeSeconds, TimeUnit.SECONDS);

        while (true) {
            try (Socket clientSocket = serverSocket.accept()) {
                BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));
                PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);

                String message = in.readLine();
                if (message.startsWith(VERIFICATION_PREFIX)) {
                    // Retornar informações do servidor em formato JSON
                    String jsonResponse = serverInfo.toJSON();
                    out.println(jsonResponse);
                    System.out.println("Resposta de /health enviada: " + jsonResponse);
                } else if ("hello".equalsIgnoreCase(message)) {
                    out.println("world");
                    System.out.println("Mensagem recebida: " + message);
                }
            } catch (IOException e) {
                System.out.println("Erro de conexão no Servidor 1: " + e.getMessage());
            }
        }
    }

    private static void sendServerInfo() throws IOException {
        try (Socket verificationSocket = new Socket(IP, VERIFICATION_PORT);
             PrintWriter verificationOut = new PrintWriter(verificationSocket.getOutputStream(), true);
             BufferedReader verificationIn = new BufferedReader(new InputStreamReader(verificationSocket.getInputStream()))) {

            // Serializa os dados do servidor para JSON
            String jsonResponse = serverInfo.toJSON();

            // Envia a solicitação de verificação
            verificationOut.println(VERIFICATION_PREFIX + "/" + jsonResponse);

            // Lê a resposta do servidor de verificação
            String verificationResponse = verificationIn.readLine();
            System.out.println("Resposta do servidor de verificação: " + verificationResponse);
        }
    }
}

