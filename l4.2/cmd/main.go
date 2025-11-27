package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AugustSerenity/go-contest-L4/l4.2_CLI-sort/internal/service"
)

func main() {
	mode := flag.String("mode", "server", "server|coord")
	addr := flag.String("addr", "", "address to listen on (server only)")
	peers := flag.String("peers", "", "comma-separated peers host:port")
	fields := flag.String("f", "", "fields")
	delimiter := flag.String("d", "\t", "delimiter")
	separated := flag.Bool("s", false, "only lines with delimiter")
	flag.Parse()

	opts := service.Options{
		Fields:    *fields,
		Delimiter: *delimiter,
		Separated: *separated,
	}

	svc := service.New(*addr, strings.Split(*peers, ","), opts)

	if *mode == "server" {
		fmt.Println("Starting server on", *addr)
		svc.RunServer(os.Stdin, os.Stdout)
		return
	}

	if *mode == "coord" {
		fmt.Fprintln(os.Stderr, "Running coordinator on", *addr)
		svc.RunCoordinator(os.Stdout)
		return
	}

	fmt.Println("unknown mode")
}
