package sthaimq

import (
	"encoding/json"
	"errors"
	"github.com/bhaswanth88/sthaimq-client-go/components"
	"github.com/bhaswanth88/sthaimq-client-go/constants"
	"github.com/bhaswanth88/sthaimq-client-go/objects"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strconv"
)

type Client struct {
	callbackChannel   chan *objects.MQControlChannel
	connectionOptions *objects.MQConnectionOptions
	nodeManager       *components.NodeManager
	subscribedTopics  map[string]bool
	wsConn            *websocket.Conn
}

func NewClient(controlChannel chan *objects.MQControlChannel) *Client {
	return &Client{callbackChannel: controlChannel}
}

func (c *Client) Connect(options *objects.MQConnectionOptions) error {
	c.connectionOptions = options
	c.nodeManager = components.NewNodeManager(options)
	broker, err := c.nodeManager.GetRandomLiveNode()
	if err != nil {
		log.Fatal("get live node:", err)
		return err
	}
	log.Println("Connecting to Live Broker:: " + broker.NodeClusterIP + ":" + strconv.Itoa(broker.NodeBrokerPort))

	wsUrl, err := url.Parse("ws://" + broker.NodeClusterIP + ":" + strconv.Itoa(broker.NodeBrokerPort))
	if err != nil {
		log.Fatal("url parse:", err)

		return err
	}

	log.Println("URL: " + wsUrl.String())
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	c.wsConn = conn
	log.Println("Connection Dialed..")
	connectEvent := new(objects.MQControlChannel)
	connectEvent.MessageType = 1
	c.callbackChannel <- connectEvent
	log.Println("Published Connect Event to Control Channel")
	c.authenticate()
	go c.readMessages()
	return nil
}

func (c *Client) readMessages() {
	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			//FIXME - if connection error, reconnect has to be done
			log.Println("read:", err)
			closeEvent := new(objects.MQControlChannel)
			closeEvent.MessageType = 3

			c.callbackChannel <- closeEvent
			log.Println("Published Close Event to Control Channel")

			return
		}
		log.Printf("recv: %s", message)

		mqMessage := new(objects.MQMessage)
		err2 := json.Unmarshal(message, mqMessage)
		if err2 != nil {
			errEvent := new(objects.MQControlChannel)
			errEvent.MessageType = 2
			c.callbackChannel <- errEvent
			log.Println("Published Err Event to Control Channel")

		} else {
			dataEvent := new(objects.MQControlChannel)
			dataEvent.MessageType = 5
			dataEvent.Body = mqMessage
			c.callbackChannel <- dataEvent
			log.Println("Published Data Event to Data Channel")

		}
	}
}

func (c *Client) Shutdown() {
	c.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.wsConn.Close()
}

func (c *Client) reSubscribeToAllTopics() {
	if c.subscribedTopics != nil {
		for topic, subscribedStatus := range c.subscribedTopics {
			if subscribedStatus {
				c.Subscribe(topic)
			}
		}
	}
}
func (c *Client) authenticate() {
	message := new(objects.MQMessage)
	message.MsgType = constants.BROKER_MSG_AUTHENTICATE
	if c.connectionOptions.GetClientID() != nil {
		message.ClientId = *c.connectionOptions.GetClientID()
		log.Println("ClientID: " + *c.connectionOptions.GetClientID())

	}

	payload := make(map[string]string)

	if c.connectionOptions.MsCliId() != nil {
		payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_MASTERCLID] = *c.connectionOptions.MsCliId()
		log.Println("MSClid: " + *c.connectionOptions.MsCliId())

	}
	if c.connectionOptions.MsToken() != nil {
		payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_MASTERTOKEN] = *c.connectionOptions.MsToken()
		log.Println("MStoken: " + *c.connectionOptions.MsToken())

	}
	if c.connectionOptions.JwtToken() != nil {
		payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_XTOKEN] = *c.connectionOptions.JwtToken()
		log.Println("JwtToken: " + *c.connectionOptions.JwtToken())

	}
	if c.connectionOptions.DeviceId() != nil {
		payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_XDEVICEID] = *c.connectionOptions.DeviceId()
		log.Println("DeviceId: " + *c.connectionOptions.DeviceId())

	}
	if c.connectionOptions.UserId() != nil {
		payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_XUID] = *c.connectionOptions.UserId()
		log.Println("UserId: " + *c.connectionOptions.UserId())

	}

	message.Payload = payload
	message.MessageId = uuid.New().String()
	err := c.sendMessage(message)
	if err != nil {
		log.Println("Failed to Send Authenticate")

	} else {
		log.Println("Authenticate Sent")
	}
	authEvent := new(objects.MQControlChannel)
	authEvent.MessageType = 4
	c.callbackChannel <- authEvent
	log.Println("Published Auth Event to Control Channel")

	c.reSubscribeToAllTopics()
}
func (c *Client) Subscribe(topic string) {
	message := new(objects.MQMessage)
	message.MsgType = constants.BROKER_MSG_SUBSCRIBE
	if c.connectionOptions.GetClientID() != nil {
		message.ClientId = *c.connectionOptions.GetClientID()
		log.Println("ClientID: " + *c.connectionOptions.GetClientID())

	}
	payload := make(map[string]string)
	payload[constants.BROKER_MSG_SUBSCRIBE_PAYLOAD_TOPIC] = topic
	message.Payload = payload
	message.MessageId = uuid.New().String()
	err := c.sendMessage(message)
	if c.subscribedTopics == nil {
		c.subscribedTopics = make(map[string]bool)
	}
	c.subscribedTopics[topic] = true
	if err != nil {
		log.Println("Failed to Subscribed For Topic: " + topic)

	} else {
		log.Println("Subscribe Sent For Topic: " + topic)
	}

}

func (c *Client) Publish(topic string, messageString string) {
	message := new(objects.MQMessage)
	message.MsgType = constants.BROKER_MSG_PUBLISH
	if c.connectionOptions.GetClientID() != nil {
		message.ClientId = *c.connectionOptions.GetClientID()
		log.Println("ClientID: " + *c.connectionOptions.GetClientID())

	}
	payload := make(map[string]string)
	payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_TOPIC] = topic
	payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_MSG] = messageString
	message.Payload = payload
	message.MessageId = uuid.New().String()
	err := c.sendMessage(message)
	if err != nil {
		log.Println("Failed to Publish For Topic: " + topic + " with message: " + messageString)

	} else {
		log.Println("Publish Sent For Topic: " + topic + " with message: " + messageString)
	}

}
func (c *Client) sendMessage(message *objects.MQMessage) error {
	if c.wsConn != nil {
		log.Println("Connection Exists, sending the message of Type: " + message.MsgType)

		bytesData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		log.Println("Message Data: " + string(bytesData))

		err2 := c.wsConn.WriteMessage(websocket.BinaryMessage, bytesData)
		if err2 != nil {
			log.Println("Exception in sending message", err2.Error())
			return err2
		}
	} else {
		log.Println("Conn is nil, cannot send message")
		return errors.New("conn: not started")
	}
	return nil
}
