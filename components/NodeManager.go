package components

import (
	"encoding/json"
	"errors"
	"github.com/bhaswanth88/sthaimq-client-go/constants"
	"github.com/bhaswanth88/sthaimq-client-go/objects"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type NodeManager struct {
	connectionOptions *objects.MQConnectionOptions
}

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
}
func NewNodeManager(connectionOptions *objects.MQConnectionOptions) *NodeManager {
	return &NodeManager{connectionOptions: connectionOptions}
}

func (n *NodeManager) GetLiveBrokers() ([]objects.BrokerNode, error) {
	log.Println("Getting Live Brokers From URL: " + *n.connectionOptions.ConnectionUrl() + constants.API_GET_GET_LIVE_NODE_URL)
	urlObject, err := url.Parse(*n.connectionOptions.ConnectionUrl() + constants.API_GET_GET_LIVE_NODE_URL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"authKey":      {*n.connectionOptions.AuthKey()},
	}

	strResponse, err := GetAppHttpClient().GetRequest(*urlObject, headers)
	if err != nil {
		return nil, err
	}
	var brokers []objects.BrokerNode
	log.Println("Live Brokers Response:: " + strResponse)

	err = json.Unmarshal([]byte(strResponse), &brokers)
	if err != nil {
		return nil, err
	}
	log.Println("Brokers Live Count:: "+ strconv.Itoa(len(brokers)))
	return brokers, nil

}
func (n *NodeManager) GetRandomLiveNode() (*objects.BrokerNode, error) {
	liveBrokers, err := n.GetLiveBrokers()
	if err != nil {
		log.Println("[NodeManager][GetRandomLiveNode] Error in selecting live brokers")
		return nil, err
	}
	if len(liveBrokers) == 0 {
		log.Println("[NodeManager][GetRandomLiveNode] No live brokers")
		return nil, errors.New("broker: 0 live brokers")
	} else {
		selectedBroker := liveBrokers[rand.Intn(len(liveBrokers))]

		log.Println("[NodeManager][GetRandomLiveNode] Selected Broker: " + selectedBroker.NodeName())
		return &selectedBroker, nil
	}
}
