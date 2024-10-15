

package rmi;

import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;

public class Servidor {

    public static void main(String[] args) {
        try {
            FuncionarioService service = new FuncionarioServiceImpl();
            Registry registry = LocateRegistry.createRegistry(1099); // cria um registro RMI na porta 1099 pros clientes se conectarem
            registry.rebind("FuncionarioService", service); // registra o serviço com o nome "FuncionarioService"
            System.out.println("Servidor RMI está rodando...");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
