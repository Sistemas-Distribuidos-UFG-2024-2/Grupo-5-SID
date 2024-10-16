import Pyro4

@Pyro4.expose
class ClassificadorNadador:
    def classificar_nadador(self, idade):
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
    daemon = Pyro4.Daemon()                
    ns = Pyro4.locateNS()                 
    uri = daemon.register(ClassificadorNadador)  
    ns.register("classificador.nadador", uri)    

    print("Servidor RMI pronto.")
    daemon.requestLoop()                   

if __name__ == "__main__":
    main()