import struct
import logging

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
    # Leo la longitud del campo
    length_bytes = recv_all(sock, 2)
    length = struct.unpack('>H', length_bytes)[0]

    # Leo los datos del campo
    data = recv_all(sock, length)
    return data.decode('utf-8')


def recv_bet(sock):
    """
    Recibe una apuesta completa del socket.
    Retorna una tupla (agency, first_name, last_name, document, birthdate, number).
    """
    agency = recv_field(sock)
    first_name = recv_field(sock)
    last_name = recv_field(sock)
    document = recv_field(sock)
    birthdate = recv_field(sock)
    number = recv_field(sock)
    return agency, first_name, last_name, document, birthdate, number


def send_response(sock, ok):
    """Envia respuesta al cliente: 1 byte (0 = OK, 1 = error)."""
    response = b'\x00' if ok else b'\x01'
    sock.sendall(response)
