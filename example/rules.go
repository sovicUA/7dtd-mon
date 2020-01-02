package main

import (
	"flag"
	"github.com/sovicUA/sq7dtd"
	"log"
)

var host = flag.String("host", "localhost", "The game server hostname or IP address")
var port = flag.Int("port", 26900, "Listener port of game server")

func main() {

	flag.Parse()

	err := sq7dtd.QueryRules(*host, *port)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get current server time
	d, h, m, err := sq7dtd.getCurrentTime()
	if err != nil {
		log.Fatal(err)
		return errors.New("Failed to convert current server time")
	}
	log.Printf("Day: %d Time: %02d:%02d", d, h, m)

}
