import socket

def main():
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect(('localhost', 12345))

    idade = input("Digite a idade do nadador: ")
    try:
        idade_int = int(idade)
        client_socket.send(str(idade_int).encode())

        classificacao = client_socket.recv(1024).decode()
        print(f"Classificação: {classificacao}")
    except ValueError:
        print("Por favor, insira um número inteiro válido para a idade.")

    client_socket.close()

if __name__ == "__main__":
    main()