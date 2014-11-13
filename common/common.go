package common

import (
	"log"
)

type ClientMessage struct {
	Filename string
}

type ServerMessage struct {
	Success  bool
	Filesize int64
}

func CheckError(prefix string, err error) {
	if err != nil {
		log.Fatalf("%s: %s\n", prefix, err.Error())
	}
}

func PrepareLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
