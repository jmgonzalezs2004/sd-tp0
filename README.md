# TP0

Juan Manuel Gonzalez Segura  
110582  

---

## Ejercicio 2

**Objetivo:** Lograr que los cambios en archivos de configuración (`config.ini` y `config.yaml`) apliquen de inmediato sin reconstruir las imágenes de Docker.

1. Modificación de script `generar.py` para inyectar las carpetas host mediantes propiedades `volumes` hacia su root correspondiente en los contenedores.
2. Eliminación de variables de entorno duras (`LOGGING_LEVEL` y `CLI_LOG_LEVEL`) del compose forzando la lectura viva obligatoria desde los volumenes inyectados en runtime.

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