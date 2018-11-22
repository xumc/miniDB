package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/xumc/miniDB/connection"
	"github.com/xumc/miniDB/utils"
)

func main() {
	Start("localhost", 3060)
}

func Start(serverIP string, serverPort int) {
	serverAddr := fmt.Sprintf("%s:%d", serverIP, serverPort)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer conn.Close()

	// listen server return information
	go printReturnInfo(conn)

	for {
		var cmd string
		inputReader := bufio.NewReader(os.Stdin)
		cmd, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		if cmd == "exit\n" {
			break
		}

		var sql = cmd
		conn.Write([]byte(connection.ConstHeader))
		contentLen := len(sql)
		conn.Write(utils.IntToBytes(contentLen))

		_, err = conn.Write([]byte(sql))
		if err != nil {
			fmt.Printf("error %s", err)
			return
		}
	}
}

func printReturnInfo(conn net.Conn) {
	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)

	go reader(readerChannel)

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			fmt.Println(conn.RemoteAddr().String(), "disconnected")
			return
		}

		if err != nil {
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = connection.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func reader(rc chan []byte) {
	for {
		select {
		case data := <-rc:
			fmt.Printf("server return: %s \n", string(data))
		}
	}
}
