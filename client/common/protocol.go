package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

// SendBetBatch envia un batch de apuestas. Primero envia la cantidad (2 bytes),
// luego cada apuesta individualmente.
func SendBetBatch(conn net.Conn, agency string, bets []Bet) error {
	// Envio la cantidad de apuestas en el batch
	countBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(countBuf, uint16(len(bets)))
	if err := sendAll(conn, countBuf); err != nil {
		return fmt.Errorf("error enviando cantidad de apuestas: %v", err)
	}

	for _, bet := range bets {
		if err := sendBet(conn, agency, bet); err != nil {
			return err
		}
	}
	return nil
}

// RecvResponse lee la respuesta del servidor (1 byte: 0 = OK, otro = error).
func RecvResponse(conn net.Conn) (bool, error) {
	resp, err := recvAll(conn, 1)
	if err != nil {
		return false, err
	}
	return resp[0] == 0, nil
}
