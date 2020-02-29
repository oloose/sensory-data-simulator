# Sensory Data Simulator with MQTT Client
Test/PoC for MQTT communication for publishing sensory data and receiving commands for configuration changes (measurement threshold adjustments, etc.)

## Includes
* Simulator for Sensory data
* MQTT-Client for data publish and command receive from and to MQTT Broker

## Requirements to Build
* Go Installed and $GOROOT and $GOPATH setup (https://golang.org/doc/install)
* [Go Dep](https://github.com/golang/dep) (move binary found under release page to $GOPATH)(see [Goland-Dep-Help](https://golang.github.io/dep/docs/new-project.html)

## Build ##
* Build mit WSL/Linux für Raspberry: ```$env GOOS=linux GOARCH=arm go build -v main.go```

## Start with CLI ##
* in "cmd" directory of project type: ```go run main.go``` for help
* Basic startup command: ```go run main.go start```
* Full startup command example: ```go run main.go -b=192.168.0.101:1883 -r=3 -a=2 -a=3 -in=5 start```
