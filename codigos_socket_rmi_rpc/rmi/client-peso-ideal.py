import Pyro4

def main():
    ns = Pyro4.locateNS()                         
    uri = ns.lookup("pesoideal.service")           
    peso_ideal_service = Pyro4.Proxy(uri)          

    altura = float(input("Digite a altura: "))
    sexo = input("Digite o sexo (m/f): ")

    peso_ideal = peso_ideal_service.calcular_peso_ideal(altura, sexo)
    if peso_ideal is not None:
        print(f"Peso ideal: {peso_ideal:.2f}")
    else:
        print("Sexo inválido.")

if __name__ == "__main__":
    main()