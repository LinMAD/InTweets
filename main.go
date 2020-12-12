package main

import (
	"flag"

	"github.com/LinMAD/InTweets/infrastructure"
	"github.com/sirupsen/logrus"
)

func main() {
	// App configuration
	host := flag.String("http_host", "localhost", "Application host")
	port := flag.String("http_port", "8080", "Application port")
	isDebug := flag.Bool("debug_mode", true, "Debug mode enabled")
	flag.Parse()

	// TODO Try to authorize Twitter before executing HTTP Server...

	// Prepare dependencies
	log := logrus.New()
	if *isDebug {
		log.Level = logrus.DebugLevel
	}

	apiServer := infrastructure.InitServerAPI(log)
	apiServer.LoadHandlers()

	// Execute
	apiServer.Run(*host + ":" + *port)
}
