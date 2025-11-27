package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func StartListener(addr string, svc Service) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "Listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			var data [][]string
			if err := json.NewDecoder(conn).Decode(&data); err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding data: %v\n", err)
				return
			}
			fmt.Fprintf(os.Stderr, "Received %d lines\n", len(data))
			svc.HandleIncomingChunk(data)
		}(conn)
	}
}
