from xmlrpc.server import SimpleXMLRPCServer

def calcular_peso_ideal(altura, sexo):
    if sexo.lower() == 'm':
        return (72.7 * altura) - 58
    elif sexo.lower() == 'f':
        return (62.1 * altura) - 44.7
    else:
        return None

def main():
    server = SimpleXMLRPCServer(('localhost', 12345))
    print("Servidor RPC esperando por conexões...")
    server.register_function(calcular_peso_ideal, 'calcular_peso_ideal')
    server.serve_forever()

if __name__ == "__main__":
    main()