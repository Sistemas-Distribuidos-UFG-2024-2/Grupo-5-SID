package rmi;

import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;
import java.util.Scanner;

public class Cliente {

    public static void main(String[] args) {
        try {
            // conecta ao registro RMI
            Registry registry = LocateRegistry.getRegistry("localhost", 1099); // o servidor está rodando na porta 1099 no localhost
            FuncionarioService service = (FuncionarioService) registry.lookup("FuncionarioService"); // FuncionarioService é o nome do serviço registrado no servidor

            Scanner scanner = new Scanner(System.in);
            System.out.print("Nome do Funcionário: ");
            String nome = scanner.nextLine();
            System.out.print("Cargo (operador/programador): ");
            String cargo = scanner.nextLine();
            System.out.print("Salário: ");
            double salario = scanner.nextDouble();

            String salarioReajustado = service.reajustarSalario(nome, cargo, salario); // chama o método remoto
            System.out.printf("Salário reajustado de %s: %s%n", nome, salarioReajustado);

            scanner.close();
        } catch (Exception e) {
            System.err.println("Erro ao conectar ao servidor RMI: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
