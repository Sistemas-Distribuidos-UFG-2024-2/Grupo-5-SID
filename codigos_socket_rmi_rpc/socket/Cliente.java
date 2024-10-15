// ex 1

package socket;
import java.io.*;
import java.net.*;

public class Cliente {
    public static void main(String[] args) {
        String servidorHost = "localhost";
        int servidorPorta = 65432;

        try (Socket socket = new Socket(servidorHost, servidorPorta)) {
            PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
            BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()));

            BufferedReader teclado = new BufferedReader(new InputStreamReader(System.in));

            System.out.print("Nome do funcionário: ");
            String nome = teclado.readLine();

            System.out.print("Cargo do funcionário: ");
            String cargo = teclado.readLine();

            System.out.print("Salário do funcionário: ");
            String salario = teclado.readLine();

            out.println(nome + "," + cargo + "," + salario);

            String resposta = in.readLine();
            System.out.println("Resposta do servidor: " + resposta);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
