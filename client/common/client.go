package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Bet           Bet
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

// StartClientLoop Envia la apuesta al servidor y espera confirmacion
func (c *Client) StartClientLoop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		log.Infof("action: signal_handler | result: success | client_id: %v | signal: %v", c.config.ID, sig)
		c.finalize()
		os.Exit(0)
	}()

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		if err := c.createClientSocket(); err != nil {
			return
		}

		// Envio la apuesta al servidor usando el protocolo custom
		err := SendBet(c.conn, c.config.ID, c.config.Bet)
		if err != nil {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
				c.config.ID, err,
			)
			c.conn.Close()
			return
		}

		// Espero confirmacion del servidor
		ok, err := RecvResponse(c.conn)
		c.conn.Close()

		if err != nil {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
				c.config.ID, err,
			)
			return
		}

		if ok {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
				c.config.Bet.Document, c.config.Bet.Number,
			)
		} else {
			log.Errorf("action: apuesta_enviada | result: fail | client_id: %v", c.config.ID)
		}

		time.Sleep(c.config.LoopPeriod)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
