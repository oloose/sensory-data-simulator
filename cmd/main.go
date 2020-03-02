package main

import (
	"fmt"
	"github.com/oloose/mqtt"
	"github.com/oloose/raspb"
	"github.com/oloose/util"
	"github.com/urfave/cli"
	"log"
	"math/rand"
	"os"
	"time"
)

var AppSetup *appSetup

const ArduinoDefaultCount = 1

type appSetup struct {
	MQTTBrokerAddress string
	RaspberryCount    int
	ArduinoCount      cli.IntSlice
	PublishInterval   int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	util.Shutdown = []func(){} //initialize slice of shutdown functions that need to be executed before shutdown
	OpenLog()
	AppSetup = &appSetup{"127.0.0.1:1883",
		1,
		cli.IntSlice{},
		5}

	app := cli.NewApp()
	app.Name = "Semiramis MQTT Simulator"
	app.Usage = "Start and configure the simulator"
	app.Compiled = time.Now()
	//define cli flags/parameteres
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "broker, b",
			Value:       "127.0.0.1:1883",
			Usage:       "Set MQTT Broker address (e.g -b=192.168.0.101:1883)",
			Destination: &AppSetup.MQTTBrokerAddress,
		},
		cli.IntFlag{
			Name:  "raspbs, r",
			Value: 1,
			Usage: "Set number of Raspberry's to create. " +
				"Number of Raspberry's is equals to number of MQTTClients (e.g -r=1)",
			Destination: &AppSetup.RaspberryCount,
		},
		cli.IntSliceFlag{
			Name:  "ards, a",
			Value: &AppSetup.ArduinoCount,
			Usage: "Set number of Arduinos to create as list of number for each raspberry created (e.g. " +
				"-a=1 -a=2 -a=3 -> Rasp1 = 1 Ard, Rasp2 = 2 Ards, Rasp3 = 3 Ards)",
		},
		cli.IntFlag{
			Name:        "intervall, in",
			Value:       5,
			Usage:       "Set intervall to publish data in seconds",
			Destination: &AppSetup.PublishInterval,
		},
	}
	// define commands
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Starts the simulator",
			Action: Start,
		},
	}

	// Gracefull shutdown
	util.GracefulShutdown()

	// run cli
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Error during start up: '%s'", err)
	}

	<-util.GracefulStop
}

// Start app.
// Setup MQTT-Clients with raspberries and arduinos. Begin to publish sensor
// data when connected to MQTT-Broker
func Start(c *cli.Context) error {
	var err error
	i := 0
	var clients []*mqtt.Client

	// if no ard count was given set 1 as default for all raspbs
	if len(AppSetup.ArduinoCount) == 0 {
		AppSetup.ArduinoCount = append(AppSetup.ArduinoCount, 1)
	}
	fmt.Println(AppSetup)
	log.Printf("#### SETUP WITH: %v", AppSetup)

	// create mqtt clients with raspbs and ards
	i = 0
	j := 0
	for len(clients) <= AppSetup.RaspberryCount-1 {
		// if a ard count for this raspberry was given use it, else use default ard count
		if len(AppSetup.ArduinoCount) > i {
			j = AppSetup.ArduinoCount[i]
			i++
		} else {
			j = ArduinoDefaultCount
		}

		//create and append new mqtt client
		nClient := mqtt.NewClient(
			raspb.NewRaspberry(j),
			fmt.Sprintf("MQTTClient_R%d", i),
			AppSetup.MQTTBrokerAddress)
		clients = append(clients, nClient)

		//connect client to mqtt broker
		if err := nClient.Connect(); err != nil {
			log.Printf(err.Error())
			panic(err)
		}
		log.Printf("---------------------------------------------------------------")
	}

	// publish sensor data to mqtt broker and print sensor data to console
	i = 0
	for {
		for _, c := range clients {
			// publish new sensor data
			c.PublishSensorData()

			// print sensor data to console
			for ai, a := range c.Raspb().Ards() {
				fmt.Printf("R%d/Rig%d/SensorData		%v\n", c.Raspb().Id(), ai, a.SensorData())
			}
		}

		t := time.Now()
		fmt.Println("TIME: ", t.Format("02.01.06; 15:04:05"))
		fmt.Println("---------------------------------------------------------------")
		time.Sleep(time.Duration(AppSetup.PublishInterval) * time.Second)
	}

	return err
}

//Opens a log file used to print log messages to
func OpenLog() {
	file, err := os.OpenFile("sim.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	util.AddShutdown(func() {
		log.Printf("############# CLOSING LOG\n\n")
		//Log file close should be the last shutdown routine to be executed
		//to prevent log file from closing before all logs have been writen
		if err := file.Close(); err != nil {
			log.Panicln("############# ERROR CLOSING LOG")
		}
	})

	log.SetOutput(file)
	log.Println("############# START LOG")
}
