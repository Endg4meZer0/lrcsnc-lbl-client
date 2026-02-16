package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	readConfig()
	initClient()
	initOutput()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
