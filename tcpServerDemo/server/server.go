package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("start the server...")
	//创建listener
	listenner, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return
	}
	//监听并接受来自客户端的链接
	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("Error accetping", err.Error())
			return
		}
		go doServerStuff(conn)
	}

}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return
		}
		fmt.Printf("Received data: %v", string(buf[:len]))
	}
}
