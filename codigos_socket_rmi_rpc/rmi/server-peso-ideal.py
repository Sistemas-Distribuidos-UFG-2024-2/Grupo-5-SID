import Pyro4

@Pyro4.expose
class PesoIdeal:
    def calcular_peso_ideal(self, altura, sexo):
        if sexo.lower() == 'm':
            return (72.7 * altura) - 58
        elif sexo.lower() == 'f':
            return (62.1 * altura) - 44.7
        else:
            return None

def main():
    daemon = Pyro4.Daemon()                
    ns = Pyro4.locateNS()                  
    uri = daemon.register(PesoIdeal)       
    ns.register("pesoideal.service", uri)  

    print("Servidor RMI pronto.")
    daemon.requestLoop()                   

if __name__ == "__main__":
    main()