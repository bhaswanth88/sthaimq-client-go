package components

import (
	"encoding/json"
	"errors"
	"gitlab.com/instasafesdp/go-sthaimq-client/constants"
	"gitlab.com/instasafesdp/go-sthaimq-client/objects"
	"log"
	"math/rand"
	"net/url"
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
	err = json.Unmarshal([]byte(strResponse), &brokers)
	if err != nil {
		return nil, err
	}
	return brokers, nil

}
func (n *NodeManager) GetRandomLiveNode() (*objects.BrokerNode, error) {
	liveBrokers, err := n.GetLiveBrokers()
	if err != nil {
		return nil, err
	}
	if len(liveBrokers) == 0 {
		return nil, errors.New("broker: 0 live brokers")
	} else {
		selectedBroker := liveBrokers[rand.Intn(len(liveBrokers))]
		return &selectedBroker, nil
	}
}
