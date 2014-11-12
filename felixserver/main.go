package main

import (
	"fmt"
	"github.com/borncrusader/felix/common"
	"net"
	"os"
)

/*
1. Start server and listen on a given port
2. Get request from client and handle it
3. Get the file name in the request and check whether the cache already has it
4. If yes, serve it from the cache
5. If not, add it to the cache and serve it

Some caveats
1. Don't cache the file if it's more than 64MB.
2. Cache size will be 64MB as well.
3. Have an eviction policy, possibly LRU? Can you make it more generic?
*/

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: server [ip:port] [directory]")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	common.CheckError("cannot resolve address", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	common.CheckError("cannot listen", err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		conn.Write([]byte("hello world!"))
	}

	os.Exit(0)
}
