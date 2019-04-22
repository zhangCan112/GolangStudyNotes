package main

import (
	"flag"
	"fmt"
	"net"
	"syscall"
)

const maxRead = 25

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		panic("usage: host port")
	}
	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))
	listener := initServer(hostAndPort)
	for {
		conn, err := listener.Accept()
		checkError(err, "Accept:")
		go connectionHandler(conn)
	}
}

func initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	checkError(err, "Resolving address:port faild:'"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	checkError(err, "ListenerTCP:")
	println("Listening to:", listener.Addr().String())
	return listener
}

func connectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	println("Connection from:", connFrom)
	sayHello(conn)
	for {
		var ibuf = make([]byte, maxRead+1)
		length, err := conn.Read(ibuf[0:maxRead])
		ibuf[maxRead] = 0
		switch err {
		case nil:
			handleMsg(length, err, ibuf)
		case syscall.EAGAIN:
			continue
		default:
			goto DISCONNECT

		}
	}
DISCONNECT:
	err := conn.Close()
	println("Closed connection: ", connFrom)
	checkError(err, "Close: ")
}

func sayHello(to net.Conn) {
	oBuf := []byte{'L', 'e', 't', '\'', 's', ' ', 'G', 'O', '!', '\n'}
	wrote, err := to.Write(oBuf)
	checkError(err, "write: wrote"+string(wrote)+"bytes.")
}

func handleMsg(length int, err error, msg []byte) {
	if length > 0 {
		print("<", length, ":")
		var str = make([]byte, len(msg)-1)
		for i := 0; ; i++ {
			if msg[i] == 0 {
				break
			}
			str = append(str, msg[i])
		}
		print(string(str))
		print(">")
	}
}

func checkError(error error, info string) {
	if error != nil {
		panic("ERROR: " + info + " " + error.Error())
	}
}
