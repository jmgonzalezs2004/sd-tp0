import struct


# Tipos de mensaje del protocolo
MSG_BET_BATCH = 0x01
MSG_NOTIFY = 0x02
MSG_QUERY_WINNERS = 0x03


def recv_all(sock, n):
    """Lee exactamente n bytes del socket, manejando short reads."""
    data = b''
    while len(data) < n:
        chunk = sock.recv(n - len(data))
        if not chunk:
            raise OSError("Conexion cerrada mientras se leia del socket")
        data += chunk
    return data


def recv_field(sock):
    """Lee un campo del protocolo: [2 bytes longitud big-endian][datos]."""
    length_bytes = recv_all(sock, 2)
    length = struct.unpack('>H', length_bytes)[0]
    data = recv_all(sock, length)
    return data.decode('utf-8')


def recv_msg_type(sock):
    """Lee el tipo de mensaje (1 byte)."""
    data = recv_all(sock, 1)
    return data[0]


def recv_bet(sock):
    """Recibe una apuesta individual del socket."""
    agency = recv_field(sock)
    first_name = recv_field(sock)
    last_name = recv_field(sock)
    document = recv_field(sock)
    birthdate = recv_field(sock)
    number = recv_field(sock)
    return agency, first_name, last_name, document, birthdate, number


def recv_bet_batch(sock):
    """Recibe un batch de apuestas (sin el byte de tipo, ya leido)."""
    count_bytes = recv_all(sock, 2)
    count = struct.unpack('>H', count_bytes)[0]
    bets = []
    for _ in range(count):
        bet = recv_bet(sock)
        bets.append(bet)
    return bets


def send_response(sock, ok):
    """Envia respuesta al cliente: 1 byte (0 = OK, 1 = error)."""
    response = b'\x00' if ok else b'\x01'
    sock.sendall(response)


def send_winners(sock, winners):
    """
    Envia la lista de ganadores al cliente.
    Formato: [1 byte OK][2 bytes cantidad][por cada uno: 2 bytes longitud + DNI]
    """
    sock.sendall(b'\x00')  # Status OK
    count_bytes = struct.pack('>H', len(winners))
    sock.sendall(count_bytes)
    for dni in winners:
        data = dni.encode('utf-8')
        length = struct.pack('>H', len(data))
        sock.sendall(length + data)


def send_error(sock):
    """Envia respuesta de error (sorteo no listo)."""
    sock.sendall(b'\x01')
