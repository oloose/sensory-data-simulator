package util

import (
	"os"
	"os/signal"
	"syscall"
	"log"
	"math/rand"
)

//slice of shutdown functions that need to be executed before shutdown
var Shutdown []func()
var GracefulStop chan os.Signal

func AddShutdown(f func()) {
	Shutdown = append(Shutdown, f)
}

func GracefulShutdown() {
	GracefulStop = make(chan os.Signal, 1)
	signal.Notify(GracefulStop, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-GracefulStop

		log.Printf("#### SHUTING DOWN...\n")
		//execute shutdown functions
		for i := len(Shutdown) - 1; i >= 0; i-- {
			if i == 0 {
				log.Printf("#### SHUTDOWN\n")
			}
			Shutdown[i]()
		}
		log.Println("### SHUTDOWN FINAL")

		os.Exit(0)
	}()
}

func randomBool() bool {
	r := rand.Intn(3) - 1
	if r == 0 {
		// if 0 set random between -1 and 0
		r = rand.Intn(2) - 1
		//r still 0? than set 1
		if r == 0 {
			r = 1
		}
	}
	if r == 1 {
		return true
	} else {
		return false
	}
}
