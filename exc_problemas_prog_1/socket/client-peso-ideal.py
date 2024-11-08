import socket

def main():
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect(('localhost', 12345))

    altura = input("Digite a altura: ")
    sexo = input("Digite o sexo (m/f): ")

    client_socket.send(f"{altura},{sexo}".encode())

    resultado = client_socket.recv(1024).decode()
    print(resultado)

    client_socket.close()

if __name__ == "__main__":
    main()