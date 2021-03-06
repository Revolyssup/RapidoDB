package transport

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

// TranslationDriver interface demands an object which
// has an Operate method, it takes in string as input
// and returns nothing
type TranslationDriver interface {
	Operate(cmd string) (string, error)
}

// Client represents an active TCP client communicating
// with the server
type Client struct {
	conn   net.Conn
	log    *log.Logger
	driver TranslationDriver
}

// New returns a new client instance
func New(conn net.Conn, l *log.Logger, d TranslationDriver) *Client {
	c := &Client{conn, l, d}

	// Send the message to the client
	c.Msg("Successfully connected to RapidoDB. Please run AUTH <user> <pass> to access the DB")

	return c
}

// InitRead reads the input of the TCP clients and passes on the received command to the driver
// after trimming the received command
func (c *Client) InitRead() {
	for {
		// Read data from TCP client and parse it
		cmd, err := bufio.NewReader(c.conn).ReadString('\n')

		// Check for errors
		if err != nil {
			// If error is io.EOF then it indicates that the client has
			// disconnected and hence closing the connection here
			if err == io.EOF {
				c.log.Printf("Client %s disconnected", c.conn.RemoteAddr().String())
				c.conn.Close()
				return
			}

			// Log the error
			c.log.Printf("Error from client %s: %v", c.conn.RemoteAddr().String(), err)
			return
		}

		// Trim the data
		cmd = strings.TrimSpace(cmd)

		// Pass the command to the driver
		res, err := c.driver.Operate(cmd)
		if err != nil {
			c.Err(err)
			continue
		}
		c.Msg(res)
	}
}

// Msg sends a message to the client
func (c *Client) Msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

// Err sends an error message to the client
func (c *Client) Err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}
