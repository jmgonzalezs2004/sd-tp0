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

## Ejercicio 6

**Objetivo:** Envío de apuestas por lotes (batches) desde archivos CSV.

**Cambios realizados:**
- El cliente ahora lee las apuestas de un archivo `.data/agency-{ID}.csv` montado por volumen, en lugar de recibirlas por variables de entorno.
- Las apuestas se agrupan en batches de tamaño configurable (`batch.maxAmount` en `config.yaml`) y se envían en una sola conexión por batch.
- El protocolo incorpora un header de 2 bytes indicando la cantidad de apuestas en el batch, seguido de cada apuesta serializada.
- El servidor recibe el batch completo, persiste todas las apuestas y loguea `apuesta_recibida` con la cantidad procesada.
- Se agregaron reintentos de conexión en el cliente para manejar la race condition de arranque.

---
## Ejercicio 7

**Objetivo:** Notificación de fin de envío, sorteo y consulta de ganadores.

**Cambios realizados:**
- El protocolo ahora usa un byte de tipo de mensaje (`0x01` = batch, `0x02` = notificación, `0x03` = consulta de ganadores) para distinguir las distintas operaciones.
- El cliente ejecuta un flujo de 3 fases: (1) envío de batches, (2) notificación de que terminó, (3) consulta de ganadores con reintentos hasta que el sorteo esté listo.
- El servidor lleva un contador de agencias notificadas. Cuando todas las agencias (N configurable via `CANT_AGENCIAS`) notifican, realiza el sorteo usando `load_bets()` y `has_won()`, y agrupa los ganadores por agencia.
- La respuesta de ganadores solo se envía después del sorteo. Antes de eso, el servidor responde con error y el cliente reintenta.

---
## Ejercicio 8

**Objetivo:** Procesamiento concurrente de conexiones en el servidor.

**Cambios realizados:**
- El servidor ahora lanza un `threading.Thread` por cada conexión entrante, permitiendo atender múltiples clientes en paralelo.
- Se utiliza `threading.Lock` para sincronizar el acceso a `store_bets()` y al estado compartido (set de agencias notificadas, resultados del sorteo), dado que `store_bets` no es thread-safe.
- Se utiliza `threading.Event` para señalizar la finalización del sorteo a los threads que consultan ganadores.

*Nota:* Se eligió multithreading sobre multiprocessing porque el GIL de Python no afecta significativamente a este caso de uso donde el cuello de botella es I/O (red y disco), no CPU.

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