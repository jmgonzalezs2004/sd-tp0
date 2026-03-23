#!/bin/bash

EXPECTED_MSG="el mensaje del papu :v"

# Creo un contenedor de busybox (Contiene netcat es liviano y contiene utilidades de unix), envio un mensaje al servidor y guardo la respuesta
RESPONSE=$(docker run --rm --network tp0_testing_net busybox sh -c "echo '$EXPECTED_MSG' | nc server 12345")

if [ "$RESPONSE" == "$EXPECTED_MSG" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
