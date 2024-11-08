import Pyro4

def main():
    ns = Pyro4.locateNS()                          
    uri = ns.lookup("classificador.nadador")      
    classificador = Pyro4.Proxy(uri)               

    idade = input("Digite a idade do nadador: ")
    try:
        idade_int = int(idade)
        classificacao = classificador.classificar_nadador(idade_int)
        print(f"Classificação: {classificacao}")
    except ValueError:
        print("Por favor, insira um número inteiro válido para a idade.")

if __name__ == "__main__":
    main()