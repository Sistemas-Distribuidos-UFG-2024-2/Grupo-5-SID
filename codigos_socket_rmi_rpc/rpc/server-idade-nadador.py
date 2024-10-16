import socket

def classificar_nadador(idade):
    if 5 <= idade <= 7:
        return "infantil A"
    elif 8 <= idade <= 10:
        return "infantil B"
    elif 11 <= idade <= 13:
        return "juvenil A"
    elif 14 <= idade <= 17:
        return "juvenil B"
    elif idade >= 18:
        return "adulto"
    else:
        return "idade fora das categorias"

def main():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(('localhost', 12345))
    server_socket.listen(1)
    print("Servidor esperando por conexões...")

    while True:
        client_socket, addr = server_socket.accept()
        print(f"Conexão de {addr} estabelecida.")

        data = client_socket.recv(1024).decode()
        idade = int(data)

        classificacao = classificar_nadador(idade)
        client_socket.send(classificacao.encode())

        client_socket.close()

if __name__ == "__main__":
    main()