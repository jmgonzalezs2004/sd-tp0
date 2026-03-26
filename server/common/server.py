import socket
import logging
import signal

from common.utils import Bet, store_bets
from common.protocol import recv_bet, send_response


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        signal.signal(signal.SIGTERM, self.handle_signal)
        signal.signal(signal.SIGINT, self.handle_signal)

    def handle_signal(self, signum, frame):
        logging.info(f"action: signal_handler | result: success | signal: {signum}")
        self.finalize()

    def finalize(self):
        logging.info("action: finalize | result: success")
        self._running = False
        self._server_socket.close()

    def run(self):
        """
        Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError:
                if not self._running:
                    logging.info('action: server_shutdown | result: success')
                    break
                else:
                    raise

    def __handle_client_connection(self, client_sock):
        """
        Read bet from a specific client socket, store it, and close the socket.

        If a problem arises in the communication with the client, the
        client socket will also be closed.
        """
        try:
            # Recibo la apuesta usando el protocolo custom
            agency, first_name, last_name, document, birthdate, number = recv_bet(client_sock)

            # Creo el objeto Bet de la catedra y lo persisto
            bet = Bet(agency, first_name, last_name, document, birthdate, number)
            store_bets([bet])

            logging.info(f'action: apuesta_almacenada | result: success | dni: {document} | numero: {number}')

            # Confirmo al cliente
            send_response(client_sock, True)
        except OSError as e:
            logging.error(f"action: apuesta_almacenada | result: fail | error: {e}")
            try:
                send_response(client_sock, False)
            except OSError:
                pass
        except Exception as e:
            logging.error(f"action: apuesta_almacenada | result: fail | error: {e}")
            try:
                send_response(client_sock, False)
            except OSError:
                pass
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
