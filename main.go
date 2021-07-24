package main

import (
	"./internal/core/xconfig"
	"./internal/server"
	"flag"
	"log"
	"syscall"
)

var (
	address  = flag.String("address", ":8080", "Address on which to listen. (Examples: localhost:8000, 127.0.0.1:9090, default=':8080')")
	compress = flag.Bool("compress", false, "Whether to enable the response compressed or not. (default=false)")
	help     = flag.Bool("help", false, "Show the help")
)

func main() {

	xconfig.ConfigPut("User_test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")

	flag.Parse()

	if *help {
		log.Println("HELP ME!!!")
		syscall.Exit(0)
	}

	server.RunMainServer(address, compress)

	log.Println("Terminated.")
}
