package main

import (
	"fmt"
	"github.com/borncrusader/felix/common"
	"net"
	"os"
)

/*
1. Start client and talk to the server
2. Request a file given through the command line
3. Receive file from server and save it in the disk
*/

func main() {
	var buf [512]byte

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: client [ip:port] [file]")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	common.CheckError("cannot resolve address", err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError("cannot connect to server", err)

	n, err := conn.Read(buf[0:])
	common.CheckError("cannot read from server", err)

	fmt.Printf("number of characters: %d; response from server: %s\n", n, buf)

	os.Exit(0)
}
