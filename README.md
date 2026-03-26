# TP0

Juan Manuel Gonzalez Segura  
110582  

---
## Ejercicio 4

**Objetivo:** Implementación de un manejo de señales básico en el servidor para permitir un apagado controlado (Graceful Shutdown).

Se ha realizado una implementación mínima de manejo de señales (`SIGTERM` y `SIGINT`) en el servidor. Esta funcionalidad permite capturar la señal de terminación, cerrar el socket de escucha y finalizar la ejecución de manera limpia.

*Nota:* Esta es una implementación inicial que será expandida a medida que se incorpore la lógica de negocio en los siguientes ejercicios.

---
## Ejercicio 5

**Objetivo:** Implementación de la lógica de negocio de Lotería Nacional. El cliente actúa como agencia de quiniela y el servidor como central de lotería.

**Cambios realizados:**
- Se definió un protocolo de comunicación propio basado en campos de longitud prefijada (2 bytes big-endian + datos). Cada apuesta se transmite como una secuencia de 6 campos: agencia, nombre, apellido, documento, nacimiento y número.
- El cliente lee los datos de la apuesta desde variables de entorno (`NOMBRE`, `APELLIDO`, `DOCUMENTO`, `NACIMIENTO`, `NUMERO`) y los envía al servidor.
- El servidor recibe la apuesta, la persiste con `store_bets()` y responde con un byte de confirmación (0 = OK).
- Se separó la capa de comunicación (`protocol.go` / `protocol.py`) del modelo de dominio (`bet.go` / `utils.py`) en ambos componentes.
- Se maneja correctamente short read y short write en ambos lados.

**Flujo del sistema:**
1. Cliente se conecta al servidor.
2. Cliente serializa y envía la apuesta usando el protocolo.
3. Servidor recibe y deserializa la apuesta, la almacena y envía confirmación.
4. Cliente recibe la confirmación y loguea `apuesta_enviada`.
5. Servidor loguea `apuesta_almacenada`.

---

## Explicación del sistema

Para generar el compose se utiliza un script de python que recibe como argumentos el nombre del archivo de salida y la cantidad de clientes. El script genera un archivo yaml con la configuración del compose.

Para levantar el programa se puede utilizar el Makefile con:

`make docker-compose-up`

`make docker-compose-logs`

Para detener la ejecución y limpiar los recursos levantados, se corre:

`make docker-compose-down`

**Administración de Configuración y Entorno:** 
Los nodos y clientes se actualizan directamente desde el Host cambiando `config.ini` o `config.yaml`, ya que los túneles locales suplen a los archivos construidos por `Dockerfile` mediante inyección montada vía Volúmenes y no interviene el `docker build` en el ciclo de vida de los cambios pequeños.

**Testeo de la red:** Usando el script
`validar-echo-server.sh`
es posible validar el funcionamiento del servidor de manera independiente interactuando directamente a través de la red de contenedores.

**Protocolo de serialización:**  
La comunicación entre cliente y servidor utiliza un protocolo propio de longitud prefijada. Cada campo se envía como 2 bytes (big-endian) indicando la longitud, seguido de los datos. Este esquema permite transmitir campos de longitud variable sin depender de delimitadores ni librerías externas, y facilita el manejo de short read/write.

**Flujo general del sistema:**  
1. El cliente se conecta al servidor por TCP.
2. El cliente serializa los datos usando el protocolo y los envía.
3. El servidor recibe, deserializa, procesa y persiste la información.
4. El servidor responde con un byte de confirmación.
5. El cliente recibe la respuesta, loguea el resultado y cierra la conexión.