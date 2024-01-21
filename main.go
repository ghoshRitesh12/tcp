package main

import (
	"flag"
	"log"

	"bruh.com/server"
)

var (
	PORT = "3000"
)

func initFlags() {
	flag.StringVar(&PORT, "port", "3000", "server port")
	flag.Parse()
}

func main() {
	initFlags()

	app := server.NewServer(PORT)
	log.Fatalln(app.Listen())
}
