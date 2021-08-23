package objects

type BrokerNode struct {
	NodeName        string `json:"nodeName"`
	NodeClusterPort int    `json:"nodeClusterPort"`
	NodeBrokerPort  int    `json:"nodeBrokerPort"`
	NodeClusterIP   string `json:"nodeClusterIp"`
	Priority        int    `json:"priority"`
	Online          bool   `json:"online"`
	Role            int    `json:"role"`
}

