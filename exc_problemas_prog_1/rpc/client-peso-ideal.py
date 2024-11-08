import xmlrpc.client

def main():
    proxy = xmlrpc.client.ServerProxy('http://localhost:12345/')

    altura = float(input("Digite a altura: "))
    sexo = input("Digite o sexo (m/f): ")

    peso_ideal = proxy.calcular_peso_ideal(altura, sexo)
    if peso_ideal is not None:
        print(f"Peso ideal: {peso_ideal:.2f}")
    else:
        print("Sexo inválido.")

if __name__ == "__main__":
    main()