package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Tipos de mensaje del protocolo
const (
	MsgBetBatch  byte = 0x01
	MsgNotify    byte = 0x02
	MsgQueryWinners byte = 0x03
)

// sendAll escribe todos los bytes en el socket, manejando short writes.
func sendAll(conn net.Conn, data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := conn.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

// recvAll lee exactamente n bytes del socket, manejando short reads.
func recvAll(conn net.Conn, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// sendField envia un campo como [2 bytes longitud big-endian][datos].
func sendField(conn net.Conn, field string) error {
	data := []byte(field)
	length := uint16(len(data))

	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, length)
	if err := sendAll(conn, lenBuf); err != nil {
		return fmt.Errorf("error enviando longitud: %v", err)
	}
	if err := sendAll(conn, data); err != nil {
		return fmt.Errorf("error enviando datos: %v", err)
	}
	return nil
}

// sendBet serializa y envia una apuesta individual por el socket.
func sendBet(conn net.Conn, agency string, bet Bet) error {
	fields := []string{
		agency,
		bet.FirstName,
		bet.LastName,
		bet.Document,
		bet.Birthdate,
		fmt.Sprintf("%d", bet.Number),
	}
	for _, field := range fields {
		if err := sendField(conn, field); err != nil {
			return err
		}
	}
	return nil
}

// SendBetBatch envia un batch de apuestas precedido por el tipo de mensaje y la cantidad.
func SendBetBatch(conn net.Conn, agency string, bets []Bet) error {
	// Tipo de mensaje
	if err := sendAll(conn, []byte{MsgBetBatch}); err != nil {
		return err
	}

	// Cantidad de apuestas
	countBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(countBuf, uint16(len(bets)))
	if err := sendAll(conn, countBuf); err != nil {
		return err
	}

	for _, bet := range bets {
		if err := sendBet(conn, agency, bet); err != nil {
			return err
		}
	}
	return nil
}

// SendNotify envia la notificación de que la agencia terminó de enviar apuestas.
func SendNotify(conn net.Conn, agency string) error {
	if err := sendAll(conn, []byte{MsgNotify}); err != nil {
		return err
	}
	return sendField(conn, agency)
}

// SendQueryWinners envia la consulta de ganadores para una agencia.
func SendQueryWinners(conn net.Conn, agency string) error {
	if err := sendAll(conn, []byte{MsgQueryWinners}); err != nil {
		return err
	}
	return sendField(conn, agency)
}

// RecvResponse lee la respuesta del servidor (1 byte: 0 = OK, otro = error).
func RecvResponse(conn net.Conn) (bool, error) {
	resp, err := recvAll(conn, 1)
	if err != nil {
		return false, err
	}
	return resp[0] == 0, nil
}

// RecvWinners lee la lista de ganadores: [2 bytes cantidad][campo DNI por cada uno].
func RecvWinners(conn net.Conn) ([]string, error) {
	// Primero leo el byte de status
	status, err := recvAll(conn, 1)
	if err != nil {
		return nil, err
	}
	if status[0] != 0 {
		return nil, fmt.Errorf("el servidor respondio con error")
	}

	// Leo la cantidad de ganadores
	countBuf, err := recvAll(conn, 2)
	if err != nil {
		return nil, err
	}
	count := binary.BigEndian.Uint16(countBuf)

	// Leo cada DNI ganador
	winners := make([]string, 0, count)
	for i := uint16(0); i < count; i++ {
		lenBuf, err := recvAll(conn, 2)
		if err != nil {
			return nil, err
		}
		length := binary.BigEndian.Uint16(lenBuf)
		data, err := recvAll(conn, int(length))
		if err != nil {
			return nil, err
		}
		winners = append(winners, string(data))
	}
	return winners, nil
}
