package main

import (
	"encoding/json"
	"fmt"
	"github.com/borncrusader/felix/cache"
	"github.com/borncrusader/felix/common"
	"log"
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
	common.PrepareLogger()

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: server [ip:port] [directory]")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	common.CheckError(fmt.Sprintf("Cannot Resolve: '%s'", os.Args[1]), err)

	err = os.Chdir(os.Args[2])
	common.CheckError(fmt.Sprintf("Cannot Chdir to '%s'", os.Args[2]), err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	common.CheckError(fmt.Sprintf("Cannot ListenTCP at '%s'", tcpAddr.String()),
		err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		handleClient(conn)

		conn.Close()
	}

	os.Exit(0)
}

func handleClient(conn net.Conn) {
	log.Printf("Client Connected: '%s'\n", conn.RemoteAddr().String())

	var buf [1024]byte

	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Printf("Cannot Read from '%s': %s", conn.RemoteAddr().String(),
			err.Error())
		return
	}

	var req common.ClientMessage

	err = json.Unmarshal(buf[:n], &req)
	if err != nil {
		log.Printf("Cannot Unmarshal request from '%s': %s",
			conn.RemoteAddr().String(), err.Error())
		return
	}

	var rsp common.ServerMessage

	/* check for the file in the given directory */
	fi, err := os.Stat(req.Filename)
	if err != nil || fi.IsDir() {
		rsp.Success = false
	} else {
		rsp.Success = true
		rsp.Filesize = fi.Size()
	}

	/* return stat of file */
	rsp_m, err := json.Marshal(rsp)
	if err != nil {
		log.Printf("Cannot Marshall response to '%s': %s\n",
			conn.RemoteAddr().String(), err.Error())
		return
	}

	_, err = conn.Write([]byte(rsp_m))
	if err != nil {
		log.Printf("Cannot Write to '%s': %s\n", conn.RemoteAddr().String(),
			err.Error())
		return
	}

	if !rsp.Success {
		log.Printf("Cannot Serve: '%s'\n", req.Filename)
		return
	}

	/* return file */
	_, ok := cache.RetrieveFile(fi)
	if !ok {
		/* read file and stream to client */
		file, err := os.Open(req.Filename)
		if err != nil {
			log.Printf("Cannot Open: '%s': %s\n", req.Filename, err.Error())
			return
		}

		var totalSize int64 = 0

		for totalSize < fi.Size() {
			n, err = file.Read(buf[0:])
			if err != nil {
				log.Printf("Cannot Read from '%s': %s\n", req.Filename,
					err.Error())
				break
			}

			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Println("Cannot Write to '%s': %s\n",
					conn.RemoteAddr().String(), err.Error())
				break
			}

			totalSize += int64(n)
		}

		file.Close()
	}
}
