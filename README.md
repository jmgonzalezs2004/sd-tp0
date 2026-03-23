# TP0

Juan Manuel Gonzalez Segura  
110582  

---


## Ejercicio 3

**Objetivo:** Validar el funcionamiento del servidor de manera independiente interactuando directamente a través de la red de contenedores.

1. Creación de un script `validar-echo-server.sh` que levanta un contenedor temporal de la imagen `busybox`, enrutándolo a la red interna compartida (`testing_net`). 
2. Utilización de la herramienta `nc` (netcat) dentro de este contenedor para enviar un mensaje directo al contenedor `server` por el puerto `12345`.
3. El script verifica que la respuesta recibida sea idéntica al payload enviado ("el mensaje del papu :v") para corroborar la correctitud del "Echo server", imprimiendo el resultado (`success` o `fail`) por salida estándar.

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