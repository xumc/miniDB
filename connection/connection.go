package connection

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

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
		c.logger.Fatalf("Error listening:", err.Error())
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
	var session []sqlparser.SQL

	for {
		select {
		case data := <-readerChannel:
			c.logger.Printf("data: %s", string(data))
			sql := strings.Trim(string(data), "\n")
			if sql == "BEGIN" {
				session = make([]sqlparser.SQL, 0)
			} else if sql == "COMMIT" {
				c.logger.Printf("sql count in the transaction : %d", len(session))
			} else {
				parsedSQL, err := c.parser.Parse(sql)
				if err != nil {
					c.logger.Printf("parser error: %s", err)
					session = nil

					errInfo := []byte(err.Error())
					conn.Write([]byte(ConstHeader))
					contentLen := len(errInfo)
					conn.Write(utils.IntToBytes(contentLen))
					conn.Write(errInfo)

					continue
				}

				session = append(session, parsedSQL)
			}

		}
	}
}

//execute TODO remove execute
// func (c *Connection) execute(sql sqlparser.SQL) error {
// 	switch sqlStruct := sql.(type) {
// 	case *sqlparser.InsertSQL:
// 		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
// 		if err != nil {
// 			return err
// 		}
// 		record := c.parser.TransformInsert(sqlStruct, tableDesc)

// 		affectedRows, err := c.store.Insert(record.TableName, record)
// 		if err != nil {
// 			return err
// 		}

// 		c.logger.Printf("affected rows: %d", affectedRows)

// 	case *sqlparser.UpdateSQL:
// 		tableDesc, err := store.GetMetadataOf(*sqlStruct.TableName)
// 		if err != nil {
// 			return err
// 		}

// 		qt, setItems := c.parser.TransformUpdate(sqlStruct, tableDesc)

// 		affectedRows, err := c.store.Update(*sqlStruct.TableName, qt, setItems)
// 		if err != nil {
// 			return err
// 		}

// 		c.logger.Printf("affected rows: %d", affectedRows)

// 	case *sqlparser.SelectSQL:
// 		qt := c.parser.TransformSelect(sqlStruct)

// 		rs, err := c.store.Select(*sqlStruct.TableName, qt)
// 		if err != nil {
// 			return err
// 		}

// 		c.logger.Printf("rows: %v", rs)

// 	case *sqlparser.DeleteSQL:
// 		qt := c.parser.TransformDelete(sqlStruct)

// 		affectedRows, err := c.store.Delete(*sqlStruct.TableName, qt)
// 		if err != nil {
// 			return err
// 		}

// 		c.logger.Printf("affected rows: %d", affectedRows)

// 	default:
// 		return errors.New("unsupport sql type")
// 	}

// 	return nil
// }
