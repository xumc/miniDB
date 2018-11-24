package connection

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/xumc/miniDB/sqlparser"
	"github.com/xumc/miniDB/utils"
)

const (
	ConnType = "tcp"
	Addr     = ":3060"
)

type Connector interface {
	Run()
}

type Connection struct {
	logger *log.Logger
	parser *sqlparser.Parser
}

func NewConnection(logger *log.Logger, p *sqlparser.Parser) *Connection {
	return &Connection{
		logger: logger,
		parser: p,
	}
}

func (c *Connection) Run() {
	l, err := net.Listen(ConnType, Addr)
	if err != nil {
		c.logger.Fatalf("Error listening: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			c.logger.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go c.handleRequest(conn)
	}
}

// Handles incoming requests.
func (c *Connection) handleRequest(conn net.Conn) {
	defer conn.Close()

	c.logger.Println(conn.RemoteAddr(), " connected")

	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go c.reader(conn, readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			c.logger.Println(conn.RemoteAddr().String(), "disconnected")
			return
		}

		if err != nil {
			c.logger.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func (c *Connection) reader(conn net.Conn, readerChannel chan []byte) {
	var transactionID uuid.UUID

	for {
		select {
		case data := <-readerChannel:
			c.logger.Printf("data: %s", string(data))
			sql := strings.Trim(string(data), "\n")
			if sql == "BEGIN" {
				transactionID = uuid.Must(uuid.NewV4())
				_, err := c.parser.Next(transactionID, sqlparser.BeginSQL)
				if err != nil {
					c.logger.Printf("begin error: %s", err)
					writeErrorInfo(conn, err)
					continue
				}
				c.logger.Printf("BEGIN transaction %v:", transactionID)

			} else if sql == "COMMIT" {
				_, err := c.parser.Next(transactionID, sqlparser.CommitSQL)
				if err != nil {
					c.logger.Printf("begin error: %s", err)
					writeErrorInfo(conn, err)
					continue
				}
				c.logger.Printf("COMMIT transaction id : %v", transactionID)
			} else {
				parsedSQL, err := c.parser.Parse(sql)
				if err != nil {
					c.logger.Printf("parser error: %s", err)
					writeErrorInfo(conn, err)
					continue
				}

				_, err = c.parser.Next(transactionID, parsedSQL)
				if err != nil {
					c.logger.Printf("next error: %s", err)
					writeErrorInfo(conn, err)
					continue
				}
			}

		}
	}
}

func writeErrorInfo(conn net.Conn, err error) {
	errInfo := []byte(err.Error())
	conn.Write([]byte(ConstHeader))
	contentLen := len(errInfo)
	conn.Write(utils.IntToBytes(contentLen))
	conn.Write(errInfo)
}
