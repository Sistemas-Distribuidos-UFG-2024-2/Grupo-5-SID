package rmi;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;

public class FuncionarioServiceImpl extends UnicastRemoteObject implements FuncionarioService {

    public FuncionarioServiceImpl() throws RemoteException {
        super();
    }

    @Override
    public String reajustarSalario(String nome, String cargo, double salario) throws RemoteException {
        double salarioReajustado = salario;

        if (cargo.equalsIgnoreCase("operador")) {
            salarioReajustado = salario * 1.20;
        } else if (cargo.equalsIgnoreCase("programador")) {
            salarioReajustado = salario * 1.18;
        }

        return String.format("Funcionário: %s, Salário Reajustado: %.2f", nome, salarioReajustado);
    }
}
