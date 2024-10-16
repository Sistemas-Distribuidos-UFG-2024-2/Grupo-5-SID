import socket

def main():
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect(('localhost', 12345))

    idade = input("Digite a idade do nadador: ")
    client_socket.send(idade.encode())

    classificacao = client_socket.recv(1024).decode()
    print(f"Classificação: {classificacao}")

    client_socket.close()

if __name__ == "__main__":
    main()