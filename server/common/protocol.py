import struct


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


def recv_bet(sock):
    """
    Recibe una apuesta individual del socket.
    Retorna una tupla (agency, first_name, last_name, document, birthdate, number).
    """
    agency = recv_field(sock)
    first_name = recv_field(sock)
    last_name = recv_field(sock)
    document = recv_field(sock)
    birthdate = recv_field(sock)
    number = recv_field(sock)
    return agency, first_name, last_name, document, birthdate, number


def recv_bet_batch(sock):
    """
    Recibe un batch de apuestas. Primero lee la cantidad (2 bytes big-endian),
    luego recibe cada apuesta.
    Retorna una lista de tuplas.
    """
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
