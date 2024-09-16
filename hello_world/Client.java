package hello_world;

import java.io.*;
import java.net.*;

public class Client {
    public static final int LOAD_BALANCER_PORT = 5610;
    public static final String IP = "localhost";

    public static void main(String[] args) {
        int[] ports = {5611};

        while (true) {
            try (Socket socket = new Socket("localhost", 5611);
                 PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
                 BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))) {

                out.println("hello");
                String response = in.readLine();

                System.out.println("Resposta do servidor " + socket.getPort() + " : " + response);
                Thread.sleep(3000);
            } catch (IOException e) {
                System.out.println("Servidor indisponível. Tentando novamente...");
            } catch (InterruptedException e) {
                System.out.println("Erro no sleep");
            } catch (Exception e) {
                System.out.println("Erro no sleep");
            }
        }
    }
}
