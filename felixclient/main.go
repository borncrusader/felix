package main

import (
	"encoding/json"
	"fmt"
	"github.com/borncrusader/felix/common"
	"log"
	"net"
	"os"
)

/*
1. Start client and talk to the server
2. Request a file given through the command line
3. Receive file from server and save it in the disk
*/

func main() {
	common.PrepareLogger()

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: client [ip:port] [file]")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	common.CheckError(fmt.Sprintf("Cannot Resolve: '%s'", os.Args[1]), err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(fmt.Sprintf("Cannot Connect to '%s'",
		tcpAddr.String()),
		err)

	log.Printf("Connected to '%s'", tcpAddr.String())

	req := common.ClientMessage{os.Args[2]}

	req_m, err := json.Marshal(req)
	common.CheckError("Cannot Marshal request", err)

	_, err = conn.Write([]byte(req_m))
	common.CheckError(fmt.Sprintf("Cannot Write to '%s'", tcpAddr.String()),
		err)

	var buf [1024]byte

	n, err := conn.Read(buf[0:])
	common.CheckError(fmt.Sprintf("Cannot Read from '%s'", tcpAddr.String()),
		err)

	var rsp common.ServerMessage

	err = json.Unmarshal(buf[:n], &rsp)
	common.CheckError("Cannot Unmarshal response", err)

	if !rsp.Success {
		log.Fatalf("Cannot Fetch: '%s'\n", req.Filename)
	}

	ret := 0

	file, err := os.Create(os.Args[2])
	if err != nil {
		ret = 1
		log.Printf("Cannot Create: '%s': %s\n", os.Args[2], err.Error())
	} else {
		var totalSize int64 = 0

		for totalSize < rsp.Filesize {
			n, err := conn.Read(buf[0:])
			if err != nil {
				ret = 1
				log.Printf("Cannot Read from '%s'\n", err.Error())
				break
			}

			n, err = file.Write(buf[:n])
			if err != nil {
				ret = 1
				log.Printf("Cannot Write to '%s'\n", err.Error())
				break
			}

			totalSize += int64(n)
		}

		file.Close()
	}

	conn.Close()

	os.Exit(ret)
}
