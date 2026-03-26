# TP0

Juan Manuel Gonzalez Segura  
110582  

---
## Ejercicio 4

**Objetivo:** Implementación de un manejo de señales básico en el servidor para permitir un apagado controlado (Graceful Shutdown).

Se ha realizado una implementación mínima de manejo de señales (`SIGTERM` y `SIGINT`) en el servidor. Esta funcionalidad permite capturar la señal de terminación, cerrar el socket de escucha y finalizar la ejecución de manera limpia.

*Nota:* Esta es una implementación inicial que será expandida a medida que se incorpore la lógica de negocio en los siguientes ejercicios.

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