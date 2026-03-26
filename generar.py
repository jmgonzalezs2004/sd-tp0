import sys

def main():
    if len(sys.argv) < 3:
        print("Uso: python3 generar.py <archivo_salida> <cantidad_clientes>")
        sys.exit(1)
        
    outfile = sys.argv[1]
    
    try:
        clients = int(sys.argv[2])
    except ValueError:
        print("La cantidad de clientes debe ser un número entero.")
        sys.exit(1)

    yaml = f"""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - CANT_AGENCIAS={clients}
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""
    for i in range(1, clients + 1):
        yaml += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/agency-{i}.csv
    depends_on:
      - server
"""

    yaml += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

    with open(outfile, 'w') as f:
        f.write(yaml)
    
    print(f"Archivo {outfile} generado exitosamente con {clients} clientes.")

if __name__ == "__main__":
    main()
