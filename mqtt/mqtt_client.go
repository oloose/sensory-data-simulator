package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/oloose/raspb"
	"github.com/oloose/util"
	"log"
	"time"
)

type Client struct {
	raspb      *raspb.Raspberry
	clientId   string
	broker     string
	mqttClient MQTT.Client
}

// Handle lost connection
var conLostHandler MQTT.ConnectionLostHandler = func(mClient MQTT.Client, mErr error) {
	opts := mClient.OptionsReader()
	log.Printf("### CLIENT CONNECTION LOST - CLIENT_ID = %s\n", opts.ClientID())
	log.Printf("# CLIENT CONNECTION LOST - MESSAGE = %s\n", mErr.Error())
}

// Default handler for unknown subscription messages
var defMessageHandler MQTT.MessageHandler = func(mClient MQTT.Client, mMsg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", mMsg.Topic())
	fmt.Printf("MSG: %s\n", mMsg.Payload())
}

// Handle successful connection
var onConnect MQTT.OnConnectHandler = func(mClient MQTT.Client) {
	opts := mClient.OptionsReader()
	log.Printf("### CLIENT CONNECTED - CLIENT_ID = %s\n", opts.ClientID())
}

// Initializes and returns a new MQTT-Client. Sets client options and shutdown handling.
// mRasp = Raspberry governing arduinos which sensor data will be published to the MQTT-Broker
// mClientId = The MQTT-Client-Id that will be used to communicate with the MQTT-Broker
// (every Id has to be unique)
// mBroker = tcp-Address of the MQTT-Broker to publish/subscribe to
func NewClient(mRasp *raspb.Raspberry, mClientId string, mBroker string) *Client {
	//init new MQTTClient
	nClient := Client{}
	nClient.raspb = mRasp
	nClient.clientId = mClientId
	nClient.broker = mBroker

	//set client options
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", mBroker))
	opts.SetClientID(fmt.Sprint("MQTTClient_R", mRasp.Id()))
	opts.AutoReconnect = true
	opts.SetConnectionLostHandler(conLostHandler)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetDefaultPublishHandler(defMessageHandler)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetMessageChannelDepth(50)
	opts.SetOnConnectHandler(onConnect)

	// create client with above client settings
	c := MQTT.NewClient(opts)
	nClient.mqttClient = c

	// Set shutdown handling
	util.AddShutdown(func() {
		nClient.Disconnect()
	})

	return &nClient
}

// Publish sensor data from all Arduinos that this client has to the MQTT-Broker
func (rClient *Client) PublishSensorData() {
	ards := rClient.raspb.Ards()

	for i := 0; i < int(rClient.raspb.UsedArdIds()); i++ {
		ard := ards[uint8(i)]

		// create publish topic
		topic := fmt.Sprintf("/R%d/Rig%d", rClient.raspb.Id(), ard.Id())

		// use sensor data (converted from struct/string-respresentation to byte[])
		// as payload for publish
		data := []byte(fmt.Sprintf("%v", ard.SensorData()))

		// publish message to 'topic' at qos = 1 and wait for the receipt
		// from server after sending message
		token := rClient.mqttClient.Publish(topic, 0, false, data)
		token.Wait()
	}
}

// Establish client connection to MQTT-Broker
func (rClient *Client) Connect() error {
	if rClient.mqttClient.IsConnected() == false {
		if token := rClient.mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}

// Disconnect client from MQTT-Broker
func (rClient *Client) Disconnect() {
	if rClient.mqttClient.IsConnected() {
		opts := rClient.mqttClient.OptionsReader()
		log.Printf("### CLIENT DISCONNECTING ... - CLIENT_ID = %s\n",
			opts.ClientID())
		rClient.mqttClient.Disconnect(250)
	}
}

func (rClient *Client) Raspb() *raspb.Raspberry {
	rasp := *rClient.raspb
	return &rasp
}

func (rClient *Client) BrokerAddress() string {
	return rClient.broker
}

func (rClient *Client) ClientId() string {
	return rClient.clientId
}

func (rClient *Client) MQTTClient() MQTT.Client {
	return rClient.mqttClient
}
