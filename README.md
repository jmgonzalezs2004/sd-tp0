# TP0

Juan Manuel Gonzalez Segura  
110582  

---

## Ejercicio

Implemento generador de config para el docker compose basado en la interfaz pública provista por la cátedra.

La implementación se realizó en python.

**Posibles mejoras a agregar:**
- Soporte para parámetros predeterminados (ej: cantidad de clientes por defecto si no se le pasa ningún valor).
- Parámetro tipo flag (ej: `--run`) para que el script no sólo genere el archivo sino que automáticamente haga un `docker compose up` ni bien termina.

---

## Explicación del sistema

Para generar el compose se utiliza un script de python que recibe como argumentos el nombre del archivo de salida y la cantidad de clientes. El script genera un archivo yaml con la configuración del compose.

Para levantar el programa se puede utilizar el Makefile con:

`make docker-compose-up`

Para ver los logs generados por los contenedores, se utiliza:

`make docker-compose-logs`

Para detener la ejecución y limpiar los recursos levantados, se corre:

`make docker-compose-down`