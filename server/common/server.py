import socket
import logging
import signal
import os

from common.utils import Bet, store_bets, load_bets, has_won
from common.protocol import (
    recv_msg_type, recv_bet_batch, recv_field,
    send_response, send_winners, send_error,
    MSG_BET_BATCH, MSG_NOTIFY, MSG_QUERY_WINNERS,
)


class Server:
    def __init__(self, port, listen_backlog, agency_count):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._agency_count = agency_count
        self._notified_agencies = set()
        self._sorteo_done = False
        self._winners_by_agency = {}
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
        try:
            msg_type = recv_msg_type(client_sock)

            if msg_type == MSG_BET_BATCH:
                self.__handle_bet_batch(client_sock)
            elif msg_type == MSG_NOTIFY:
                self.__handle_notify(client_sock)
            elif msg_type == MSG_QUERY_WINNERS:
                self.__handle_query_winners(client_sock)
            else:
                logging.error(f"action: receive_message | result: fail | error: tipo de mensaje desconocido {msg_type}")
                send_response(client_sock, False)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __handle_bet_batch(self, client_sock):
        """Recibe y almacena un batch de apuestas."""
        raw_bets = recv_bet_batch(client_sock)
        bets = []
        for agency, first_name, last_name, document, birthdate, number in raw_bets:
            bets.append(Bet(agency, first_name, last_name, document, birthdate, number))
        store_bets(bets)
        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        send_response(client_sock, True)

    def __handle_notify(self, client_sock):
        """Registra que una agencia terminó de enviar apuestas."""
        agency = recv_field(client_sock)
        self._notified_agencies.add(agency)
        logging.info(f'action: notify | result: success | agency: {agency} | notified: {len(self._notified_agencies)}/{self._agency_count}')
        send_response(client_sock, True)

        # Si todas las agencias notificaron, realizamos el sorteo
        if len(self._notified_agencies) == self._agency_count:
            self.__do_sorteo()

    def __do_sorteo(self):
        """Realiza el sorteo: carga todas las apuestas y determina ganadores por agencia."""
        self._winners_by_agency = {}
        for bet in load_bets():
            agency_key = str(bet.agency)
            if agency_key not in self._winners_by_agency:
                self._winners_by_agency[agency_key] = []
            if has_won(bet):
                self._winners_by_agency[agency_key].append(bet.document)

        self._sorteo_done = True
        logging.info(f'action: sorteo | result: success')

    def __handle_query_winners(self, client_sock):
        """Responde con los ganadores de la agencia consultada."""
        agency = recv_field(client_sock)

        if not self._sorteo_done:
            send_error(client_sock)
            return

        winners = self._winners_by_agency.get(agency, [])
        send_winners(client_sock, winners)

    def __accept_new_connection(self):
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
