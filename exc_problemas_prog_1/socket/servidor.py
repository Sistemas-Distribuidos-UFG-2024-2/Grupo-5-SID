# ex 1
import socket


def calcular_reajuste(cargo, salario):
    if cargo.lower() == 'operador':
        return salario * 1.20
    elif cargo.lower() == 'programador':
        return salario * 1.18
    return salario


def servidor():
    servidor_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    servidor_socket.bind(('localhost', 65432))
    servidor_socket.listen(1)

    print("Servidor aguardando conexões...")

    while True:
        conn, addr = servidor_socket.accept()
        print(f"Conectado a: {addr}")

        data = conn.recv(1024).decode()
        nome, cargo, salario = data.split(",")
        salario = float(salario)

        salario_reajustado = calcular_reajuste(cargo, salario)
        resposta = f"Funcionário: {nome}, Salário Reajustado: {salario_reajustado:.2f}"
        conn.send(resposta.encode())

        conn.close()


if __name__ == "__main__":
    servidor()
