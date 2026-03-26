package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopPeriod    time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// createClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) finalize() {
	if c.conn != nil {
		c.conn.Close()
	}
	log.Infof("action: finalize | result: success | client_id: %v", c.config.ID)
}

// readBetsFromFile lee las apuestas del archivo CSV de la agencia.
// Formato esperado: Nombre,Apellido,Documento,Nacimiento,Numero
func readBetsFromFile(filepath string) ([]Bet, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo de apuestas: %v", err)
	}
	defer file.Close()

	var bets []Bet
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}
		number, _ := strconv.Atoi(fields[4])
		bet := Bet{
			FirstName: fields[0],
			LastName:  fields[1],
			Document:  fields[2],
			Birthdate: fields[3],
			Number:    number,
		}
		bets = append(bets, bet)
	}
	return bets, scanner.Err()
}

// StartClientLoop Lee las apuestas del archivo CSV y las envia al servidor en batches
func (c *Client) StartClientLoop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		log.Infof("action: signal_handler | result: success | client_id: %v | signal: %v", c.config.ID, sig)
		c.finalize()
		os.Exit(0)
	}()

	// Leo las apuestas del archivo de la agencia
	filename := fmt.Sprintf("agency-%s.csv", c.config.ID)
	bets, err := readBetsFromFile(filename)
	if err != nil {
		log.Errorf("action: read_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	// Envio las apuestas en batches
	batchSize := c.config.BatchMaxAmount
	for i := 0; i < len(bets); i += batchSize {
		end := i + batchSize
		if end > len(bets) {
			end = len(bets)
		}
		batch := bets[i:end]

		if err := c.createClientSocket(); err != nil {
			return
		}

		err := SendBetBatch(c.conn, c.config.ID, batch)
		if err != nil {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
				c.config.ID, err,
			)
			c.conn.Close()
			return
		}

		ok, err := RecvResponse(c.conn)
		c.conn.Close()

		if err != nil {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
				c.config.ID, err,
			)
			return
		}

		if !ok {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v", c.config.ID)
			return
		}

		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
