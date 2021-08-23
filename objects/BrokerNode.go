package objects

type BrokerNode struct {
	nodeName        string `json:"nodeName"`
	nodeClusterPort int    `json:"nodeClusterPort"`
	nodeClusterIp   string `json:"nodeClusterIp"`
	nodeBrokerPort  int    `json:"nodeBrokerPort"`
	isOnline        bool   `json:"online"`
}

func (b *BrokerNode) NodeName() string {
	return b.nodeName
}

func (b *BrokerNode) SetNodeName(nodeName string) {
	b.nodeName = nodeName
}

func (b *BrokerNode) NodeClusterPort() int {
	return b.nodeClusterPort
}

func (b *BrokerNode) SetNodeClusterPort(nodeClusterPort int) {
	b.nodeClusterPort = nodeClusterPort
}

func (b *BrokerNode) NodeClusterIp() string {
	return b.nodeClusterIp
}

func (b *BrokerNode) SetNodeClusterIp(nodeClusterIp string) {
	b.nodeClusterIp = nodeClusterIp
}

func (b *BrokerNode) NodeBrokerPort() int {
	return b.nodeBrokerPort
}

func (b *BrokerNode) SetNodeBrokerPort(nodeBrokerPort int) {
	b.nodeBrokerPort = nodeBrokerPort
}

func (b *BrokerNode) IsOnline() bool {
	return b.isOnline
}

func (b *BrokerNode) SetIsOnline(isOnline bool) {
	b.isOnline = isOnline
}
