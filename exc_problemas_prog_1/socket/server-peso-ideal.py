import socket

def calcular_peso_ideal(altura, sexo):
    if sexo.lower() == 'm':
        return (72.7 * altura) - 58
    elif sexo.lower() == 'f':
        return (62.1 * altura) - 44.7
    else:
        return None

def main():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(('localhost', 12345))
    server_socket.listen(1)
    print("Servidor esperando por conexões...")

    while True:
        client_socket, addr = server_socket.accept()
        print(f"Conexão de {addr} estabelecida.")

        data = client_socket.recv(1024).decode()
        altura, sexo = data.split(',')
        altura = float(altura)

        peso_ideal = calcular_peso_ideal(altura, sexo)
        if peso_ideal is not None:
            client_socket.send(f"Peso ideal: {peso_ideal:.2f}".encode())
        else:
            client_socket.send("Sexo inválido.".encode())

        client_socket.close()

if __name__ == "__main__":
    main()