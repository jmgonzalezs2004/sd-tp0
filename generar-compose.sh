#!/bin/bash
# Recibe como primer argumento el archivo de salida y como segundo la cantidad de clientes
OUTFILE="$1"
CLIENTS="$2"

echo "Nombre del archivo de salida: $OUTFILE"
echo "Cantidad de clientes: $CLIENTS"

python3 generar.py "$OUTFILE" "$CLIENTS"
