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
	ID             string
	ServerAddress  string
	LoopPeriod     time.Duration
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

// createClientSocket Initializes client socket with retries.
func (c *Client) createClientSocket() error {
	var conn net.Conn
	var err error

	for retries := 0; retries < 10; retries++ {
		conn, err = net.Dial("tcp", c.config.ServerAddress)
		if err == nil {
			c.conn = conn
			return nil
		}
		log.Debugf("action: connect | result: fail | client_id: %v | retry: %v | error: %v",
			c.config.ID, retries+1, err,
		)
		time.Sleep(500 * time.Millisecond)
	}

	log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
	return err
}

func (c *Client) finalize() {
	if c.conn != nil {
		c.conn.Close()
	}
	log.Infof("action: finalize | result: success | client_id: %v", c.config.ID)
}

// readBetsFromFile lee las apuestas del archivo CSV de la agencia.
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

// StartClientLoop Ejecuta el flujo completo: enviar apuestas, notificar, consultar ganadores
func (c *Client) StartClientLoop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		log.Infof("action: signal_handler | result: success | client_id: %v | signal: %v", c.config.ID, sig)
		c.finalize()
		os.Exit(0)
	}()

	// Fase 1: Enviar todas las apuestas en batches
	filename := fmt.Sprintf("agency-%s.csv", c.config.ID)
	bets, err := readBetsFromFile(filename)
	if err != nil {
		log.Errorf("action: read_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

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
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.conn.Close()
			return
		}

		ok, err := RecvResponse(c.conn)
		c.conn.Close()

		if err != nil || !ok {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v", c.config.ID)
			return
		}

		time.Sleep(c.config.LoopPeriod)
	}

	// Fase 2: Notificar al servidor que terminamos de enviar
	if err := c.createClientSocket(); err != nil {
		return
	}
	err = SendNotify(c.conn, c.config.ID)
	if err != nil {
		log.Errorf("action: notify | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.conn.Close()
		return
	}
	ok, err := RecvResponse(c.conn)
	c.conn.Close()
	if err != nil || !ok {
		log.Errorf("action: notify | result: fail | client_id: %v", c.config.ID)
		return
	}

	// Fase 3: Consultar ganadores (con reintentos, el sorteo puede no estar listo)
	for {
		if err := c.createClientSocket(); err != nil {
			return
		}
		err = SendQueryWinners(c.conn, c.config.ID)
		if err != nil {
			log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.conn.Close()
			return
		}
		winners, err := RecvWinners(c.conn)
		c.conn.Close()

		if err != nil {
			// El sorteo todavía no se realizó, reintento
			time.Sleep(1 * time.Second)
			continue
		}

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))
		break
	}
}
