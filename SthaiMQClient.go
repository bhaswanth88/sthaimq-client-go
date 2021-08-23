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
	controlChannel    chan *objects.MQControlChannel
	dataChannel       chan *objects.MQDataChannel
	connectionOptions *objects.MQConnectionOptions
	nodeManager       *components.NodeManager
	subscribedTopics  map[string]bool
	wsConn            *websocket.Conn
}

func NewClient(controlChannel chan *objects.MQControlChannel, dataChannel chan *objects.MQDataChannel) *Client {
	return &Client{controlChannel: controlChannel, dataChannel: dataChannel}
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
			closeEvent.SetMessageType(3)
			c.controlChannel <- closeEvent
			return
		}
		log.Printf("recv: %s", message)

		mqMessage := new(objects.MQMessage)
		err2 := json.Unmarshal(message, mqMessage)
		if err2 != nil {
			errEvent := new(objects.MQControlChannel)
			errEvent.SetMessageType(2)
			c.controlChannel <- errEvent
		} else {
			dataEvent := new(objects.MQDataChannel)
			dataEvent.SetBody(mqMessage)
			c.dataChannel <- dataEvent
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
	message.SetMsgType(constants.BROKER_MSG_AUTHENTICATE)
	message.SetClientId(*c.connectionOptions.ClientId())
	log.Println("ClientID: " + *c.connectionOptions.ClientId())

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

	message.SetPayload(payload)
	message.SetMessageId(uuid.New().String())
	err := c.sendMessage(message)
	if err != nil {
		log.Println("Failed to Send Authenticate")

	} else {
		log.Println("Authenticate Sent")
	}

	c.reSubscribeToAllTopics()
}
func (c *Client) Subscribe(topic string) {
	message := new(objects.MQMessage)
	message.SetMsgType(constants.BROKER_MSG_SUBSCRIBE)
	message.SetClientId(*c.connectionOptions.ClientId())
	payload := make(map[string]string)
	payload[constants.BROKER_MSG_SUBSCRIBE_PAYLOAD_TOPIC] = topic
	message.SetPayload(payload)
	message.SetMessageId(uuid.New().String())
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
	message.SetMsgType(constants.BROKER_MSG_PUBLISH)
	message.SetClientId(*c.connectionOptions.ClientId())
	payload := make(map[string]string)
	payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_TOPIC] = topic
	payload[constants.BROKER_MSG_PUBLISH_PAYLOAD_MSG] = messageString
	message.SetPayload(payload)
	message.SetMessageId(uuid.New().String())
	err := c.sendMessage(message)
	if err != nil {
		log.Println("Failed to Publish For Topic: " + topic + " with message: " + messageString)

	} else {
		log.Println("Publish Sent For Topic: " + topic + " with message: " + messageString)
	}

}
func (c *Client) sendMessage(message *objects.MQMessage) error {
	if c.wsConn != nil {
		bytesData, err := json.Marshal(message)
		if err != nil {
			return err
		}
		err2 := c.wsConn.WriteMessage(websocket.BinaryMessage, bytesData)
		if err2 != nil {
			return err2
		}
	} else {
		return errors.New("conn: not started")
	}
	return nil
}
