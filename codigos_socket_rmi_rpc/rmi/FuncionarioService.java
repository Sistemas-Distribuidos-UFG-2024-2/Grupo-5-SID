package rmi;

import java.rmi.Remote;
import java.rmi.RemoteException;

public interface FuncionarioService extends Remote {
    String reajustarSalario(String nome, String cargo, double salario) throws RemoteException;
}
